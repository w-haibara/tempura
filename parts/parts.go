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

func Commands(c []config.Command) string {
	str := "	if (false) {\n"
	for i, _ := range c {
		str += "	} else if(command === '" + c[i].Command + "') {\n"
		for j, _ := range c[i].Prompts {
			str += "		var var" + strconv.Itoa(j) + ";\n"
		}
		str += Message(c[i])
		str += Prompts(c[i])
	}
	str += "	} else if(command !== '') {\n"
	str += "		 term.echo(command + ': command not found')\n"
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

	pop := len(c.Prompts) >= 2

	for i, _ := range c.Prompts {
		j := len(c.Prompts) - 1 - i
		str += "		term.push(function(command, term) {\n"
		str += "			if (command) {\n"
		str += "				var" + strconv.Itoa(j) + " = command;\n"
		if i == 0 {
			str += "				var json = " + Json(c.Prompts) + ";\n"
			str += "				term.echo('send data: ' + json + '\\n');\n"
			str += Api(c.Api, pop)
		} else {
			str += "				term.pop();\n"
		}
		str += "			}\n"
		str += "		}, {\n"
		str += "			prompt: '" + c.Prompts[j].Prompt + ": '\n"
		str += "		});\n"
	}
	return str
}

func Json(p []config.Prompt) string {
	str := `'{' + `
	for i, _ := range p {
		str += `'"` + p[i].Json + `": ' + `
		str += `'"' + var` + strconv.Itoa(i) + ` + '"' + `
		if i != len(p)-1 {
			str += "',' + "
		}
	}
	str += `'}'`
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

	if a.Method == config.POST || a.Method == config.PUT {
		str += "					data: json,\n"
	}

	str += `				}).then(
					function (data) {
						console.log(data);
						term.echo('=== result ===\n' + JSON.stringify(data, null, 2));
					},
					function (jqXHR, textStatus, errorThrown) {
						term.echo("status: " + jqXHR.status + ", " + textStatus);
						term.echo(errorThrown.name + ": " + errorThrown.message);
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
