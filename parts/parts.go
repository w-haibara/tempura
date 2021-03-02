package parts

import (
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
		str += Message(c[i])
		str += Get(c[i])
	}
	return str + "	}"
}

func Message(c config.Command) string {
	return "		term.echo('" + c.Message + "');\n"
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
