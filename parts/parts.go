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

		pop := false
		str += Prompts(c[i], &pop)

		if c[i].Api != (config.Api{}) {
			str += Api(c[i].Api, pop)
		}
	}
	return str + "	}"
}

func Message(c config.Command) string {
	return "		term.echo('" + c.Message + "');\n"
}

func Prompts(c config.Command, pop *bool) string {
	str := ""

	if len(c.Prompts) >= 2 {
		*pop = true
	}

	for i, _ := range c.Prompts {
		j := len(c.Prompts) - 1 - i
		str += "		term.push(function(command, term) {\n"
		str += "			if (command) {\n"
		str += "				var" + strconv.Itoa(j) + " = command;\n"
		if i == 0 {
			str += "				var json = " + Json(c.Prompts) + ";\n"
			str += "				term.echo(json);\n"
		}
		str += "				term.pop();\n"
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

func CrossDomain() string {
	return `
(function( jQuery ) {
	if ( window.XDomainRequest ) {
		jQuery.ajaxTransport(function( s ) {
			if ( s.crossDomain && s.async ) {
				if ( s.timeout ) {
					s.xdrTimeout = s.timeout;
					delete s.timeout;
					}
				var xdr;
				return {
					send: function( _, complete ) {
						function callback( status, statusText, responses, responseHeaders ) {
							xdr.onload = xdr.onerror = xdr.ontimeout = jQuery.noop;
							xdr = undefined;
							complete( status, statusText, responses, responseHeaders );
							}
						xdr = new XDomainRequest();
						xdr.onload = function() {
							callback( 200, "OK", { text: xdr.responseText }, "Content-Type: " + xdr.contentType );
						};
						xdr.onerror = function() {
								callback( 404, "Not Found" );
						};
						xdr.onprogress = jQuery.noop;
							xdr.ontimeout = function() {
							callback( 0, "timeout" );
						};
						xdr.timeout = s.xdrTimeout || Number.MAX_VALUE;
						xdr.open( s.type, s.url );
						xdr.send( ( s.hasContent && s.data ) || null );
					},
					abort: function() {
						if ( xdr ) {
							xdr.onerror = jQuery.noop;
							xdr.abort();
						}
					}
				};
			}
		});
	}
})( jQuery );
`
}

func Api(a config.Api, pop bool) string {
	str := "		term.pause();"
	str += "		var popFlag = false;"
	str += "		$.ajax({\n"
	str += "			type: '" + a.Method + "',\n"
	str += "			url: '" + a.Url + "'\n"
	str += `
		}).then(
			function (data) {
				term.echo(JSON.stringify(data, null, 2));
			},
			function (jqXHR, textStatus, errorThrown) {
				term.echo(jqXHR.status+" : "+jqXHR.responseText);
			}
		).then(
			function () {
				term.resume();
`
	if pop {
		str += "				term.pop();\n"
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
