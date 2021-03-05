package parts

import (
	"strconv"
	"tempura/config"
)

func Header1() string {
	return `
<!DOCTYPE html>
<html>
<head>
<script src="https://code.jquery.com/jquery-3.3.1.min.js"></script>
<script src="https://unpkg.com/jquery.terminal@2.x.x/js/jquery.terminal.min.js"></script>
<link rel="stylesheet" href="https://unpkg.com/jquery.terminal@2.x.x/css/jquery.terminal.min.css"/>
`
}

func Style(strs []string) string {
	str1 := `
<style>
.terminal {
`
	str2 := `}
</style>
`
	for _, v := range strs {
		str1 += v + "\n"
	}
	return str1 + str2
}

func Header2() string {
	str := `
</head>
<body>
<script>`
	str += "$('body').terminal(function(command, term) {\n"
	return str
}

func MsgSuccess(str string) string {
	return "'[[;#8BC34A;]' + " + str + " + ']'"
}

func MsgError(str string) string {
	return "'[[;#F44336;]' + " + str + " + ']'"
}

func Commands(c []config.Command) string {
	str := "	if (false) {\n"
	for i, _ := range c {
		str += "	} else if(command === '" + c[i].Command + "') {\n"
		str += "		var history = term.history();\n"
		str += "		history.disable();\n"
		str += "		var json = ''\n"
		str += "		var header = ''\n"
		str += "		var query = ''\n"
		for j, _ := range c[i].Prompts {
			str += "		var var" + strconv.Itoa(j) + ";\n"
		}
		str += Message(c[i])
		str += Prompts(c[i])
	}
	str += "	} else if(command !== '') {\n"
	str += "		 term.echo(" + MsgError("command + ': command not found'") + ")\n"
	return str + "	}"
}

func Message(c config.Command) string {
	return "		term.echo('" + c.Message + "' + '\\n');\n"
}

func Prompts(c config.Command) string {
	if len(c.Prompts) == 0 {
		return Api(c.Api, false)
	}

	str := ""

	for i, _ := range c.Prompts {
		j := len(c.Prompts) - 1 - i
		str += "		term.push(function(command, term) {\n"
		str += "			if (command) {\n"
		if (i == 0 && c.Prompts[0].Mask) || (i < len(c.Prompts)-1 && c.Prompts[i+1].Mask) {
			str += "				term.set_mask(false);\n"
		}
		str += "				var" + strconv.Itoa(j) + " = command;\n"
		if i == 0 {
			str += InitializationOfObjects(c.Prompts)
			if c.Print.Json {
				str += "				term.echo('json: ' + JSON.stringify(json, undefined, 2));\n"
			}
			if c.Print.Header {
				str += "				term.echo('header: ' + JSON.stringify(header, undefined, 2));\n"
			}
			if c.Print.Query {
				str += "				term.echo('query: ' + JSON.stringify(query, undefined, 2));\n"
			}
			if s := Api(c.Api, true); s != "" {
				str += s
			} else {
				str += "				term.pop();\n"
			}
		} else {
			str += "				term.pop();\n"
		}
		if i != 0 && i < len(c.Prompts) && c.Prompts[i].Mask {
			str += "				term.set_mask('*');\n"
		}
		str += "			}\n"
		str += "		}, {\n"
		str += "			prompt: '" + c.Prompts[j].Prompt + ": '\n"
		str += "		});\n"
	}

	if c.Prompts[0].Mask {
		str += "		term.set_mask('*');\n"
	}

	return str
}

func InitializationOfObjects(p []config.Prompt) string {
	str := ""
	if !config.IsEmpty(p, "Json") {
		str += "				json = " + Json(p) + ";\n"
	}
	if !config.IsEmpty(p, "Header") {
		str += "				header = " + Header(p) + ";\n"
	}
	if !config.IsEmpty(p, "Query") {
		str += "				query = " + Query(p) + ";\n"
	}
	return str
}

func Json(p []config.Prompt) string {
	str := ""
	for i, _ := range p {
		str += ToObject(i, len(p), p[i].Json)
	}
	return str
}

func Header(p []config.Prompt) string {
	str := ""
	for i, _ := range p {
		str += ToObject(i, len(p), p[i].Header)
	}
	return str
}

func Query(p []config.Prompt) string {
	str := ""
	for i, _ := range p {
		str += ToObject(i, len(p), p[i].Query)
	}
	return str
}

func ToObject(i int, length int, v string) string {
	str := ""
	if i == 0 {
		str += `{`
	}
	str += `"` + v + `": `
	str += `var` + strconv.Itoa(i) + ``
	if i != length-1 {
		str += ", "
	} else {
		str += `}`
	}
	return str
}

func Api(a config.Api, pop bool) string {
	if a == (config.Api{}) {
		return ""
	}

	str := "				term.pause();\n"
	str += "				var popFlag = false;\n"

	str += "				$.ajax({\n"
	str += "					type: '" + a.Method + "',\n"
	str += "					url: '" + a.Url + "',\n"
	str += "					crossDomain: true,\n"

	if a.Method == config.POST || a.Method == config.PUT {
		str += "					data: json,\n"
	}

	str += `				}).then(
					function (data) {
						console.log(data);
						if (typeof data == 'string') {
							data = JSON.parse(data);
						}
`
	str += "						term.echo(" +
		MsgSuccess("'=== result ===\\n' + JSON.stringify(data, null, 2)") + ");\n"
	str += `					},
					function (jqXHR, textStatus, errorThrown) {
						term.echo(` +
		MsgError(`"status: " + jqXHR.status + ", " + textStatus`) +
		`);
						term.echo(` +
		MsgError(`errorThrown.name + ": " + errorThrown.message`) +
		`);
					}
				).then(
					function () {
						term.resume();
`
	if pop {
		str += "						term.pop();\n"
	}
	str += `
					}
				);
`
	return str
}

func Footer(greeting string) string {
	str1 := `
}, {
	greetings: function(cb) {
		cb("`
	str2 := `");
	}
});
</script>
</body>
</html>
`
	return str1 + greeting + str2
}
