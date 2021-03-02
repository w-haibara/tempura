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
	return `
</head>
<body>
<script>
$('body').terminal(function(command, term) {
`
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
	return str + "	}"
}

func Message(c config.Command) string {
	return "		term.echo('" + c.Message + "');\n"
}

func Prompts(c config.Command) string {
	str := ""
	for i, _ := range c.Prompts {
		j := len(c.Prompts) - 1 - i
		str += "		term.push(function(command, term) {\n"
		str += "			if (command) {\n"
		str += "				var" + strconv.Itoa(j) + " = command;\n"
		if i == 0 {
			str += "				term.echo("
			for i, _ := range c.Prompts {
				str += "'var" + strconv.Itoa(i) + ":' + var" + strconv.Itoa(i) + " + ' ' + "
			}
			str += "'');\n"
		}
		str += "				term.pop();\n"
		str += "			}\n"
		str += "		}, {\n"
		str += "			prompt: '" + c.Prompts[j].Prompt + ": '\n"
		str += "		});\n"
	}
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
