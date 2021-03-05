package parts

import (
	"strconv"
	"strings"
	"tempura/config"
)

func Header1() string {
	return `
<!DOCTYPE html>
<html>
<head>
<script src="https://code.jquery.com/jquery-3.3.1.min.js"></script>
<script src="https://cdn.jsdelivr.net/gh/jcubic/static/js/wcwidth.js"></script>
<script src="https://unpkg.com/jquery.terminal/js/jquery.terminal.js"></script>
<link href="https://unpkg.com/jquery.terminal/css/jquery.terminal.css" rel="stylesheet"/>
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

func FormattedString(str, option string) string {
	return "'[[" + option + "]' + (" + str + ").replace(/&/g, '&amp;').replace(/\\[/g, '&#91;').replace(/\\]/g, '&#93;') + ']'"
}

func MsgSuccess(str string) string {
	return FormattedString(str, ";#8BC34A;")
}

func MsgError(str string) string {
	return FormattedString(str, ";#F44336;")
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
	str += "	} else if(command.startsWith('echo')) {\n"
	str += "		term.echo(command.replace(/^echo/, '').trimStart());\n"
	str += "	} else if(command !== '') {\n"
	str += "		term.echo(" + MsgError("command + ': command not found'") + ");\n"
	str += "		console.log(" + MsgError("command + ': command not found'") + ");\n"
	return str + "	}"
}

func Message(c config.Command) string {
	return "		term.echo('" + c.Message + "');\n"
}

func Prompts(c config.Command) string {
	if len(c.Prompts) == 0 {
		return Api(nil, c.Api, false)
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
				str += "				term.echo('\\njson: ' + JSON.stringify(json, undefined, 2));\n"
			}
			if c.Print.Header {
				str += "				term.echo('\\nheader: ' + JSON.stringify(header, undefined, 2));\n"
			}
			if c.Print.Query {
				str += "				term.echo('\\nquery: ' + JSON.stringify(query, undefined, 2));\n"
			}
			if s := Api(c.Prompts, c.Api, true); s != "" {
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
	jsonStr := ""
	headerStr := ""
	queryStr := ""
	for i, _ := range p {
		jsonStr += ToObjectOneByOne(i, len(p), p[i].Json)
		headerStr += ToObjectOneByOne(i, len(p), p[i].Header)
		queryStr += ToObjectOneByOne(i, len(p), p[i].Query)
	}

	jsonStr = "{" + strings.TrimSuffix(jsonStr, ",") + "}"
	headerStr = "{" + strings.TrimSuffix(headerStr, ",") + "}"
	queryStr = "{" + strings.TrimSuffix(queryStr, ",") + "}"

	str := ""
	if !config.IsEmpty(p, "Json") {
		str += "				json = " + jsonStr + ";\n"
	}
	if !config.IsEmpty(p, "Header") {
		str += "				header = " + headerStr + ";\n"
	}
	if !config.IsEmpty(p, "Query") {
		str += "				query = " + queryStr + ";\n"
	}
	return str
}

func ToObjectOneByOne(i int, length int, v string) string {
	str := ""
	if v != "" {
		str += `"` + v + `": `
		str += `var` + strconv.Itoa(i) + `,`
	}
	return str
}

func Api(p []config.Prompt, a config.Api, pop bool) string {
	if a == (config.Api{}) {
		return ""
	}

	str := "				term.pause();\n"
	str += "				var popFlag = false;\n"

	str += "				$.ajax({\n"
	str += "					type: '" + a.Method + "',\n"
	if !config.IsEmpty(p, "Query") {
		str += "					url: '" + a.Url + "?' + jQuery.param(query)" + ",\n"
	} else {
		str += "					url: '" + a.Url + "',\n"
	}
	str += "					crossDomain: true,\n"

	if !config.IsEmpty(p, "Json") && (a.Method == config.POST || a.Method == config.PUT) {
		str += "					data: JSON.stringify(json, undefined),\n"
		str += "					dataType: 'json',\n"
		str += "					contentType: 'application/json',\n"
	}
	if !config.IsEmpty(p, "Header") {
		str += "					headers: header\n"
	}

	str += `				}).then(
					function (data) {
						if (typeof data == 'string') {
							data = JSON.parse(data);
						}
`
	str += "						term.echo(" +
		MsgSuccess("'\\n=== result ===\\n' + JSON.stringify(data, null, 2)") + ");\n"
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
