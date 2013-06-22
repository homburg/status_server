package status_server

const ServerTemplateAmber = `doctype 5
html[lang="en"]
head
	meta[charset="UTF-8"]
	title Status page
	style[type="text/css"]
		div {
			margin: 1em;
		}
	script[type="text/javascript"][src="//cdnjs.cloudflare.com/ajax/libs/jquery/2.0.1/jquery.min.js"]
	link[href="//netdna.bootstrapcdn.com/twitter-bootstrap/2.2.2/css/bootstrap-combined.min.css"][rel="stylesheet"]
body
	div
		h1 #{Hostname}
		p.muted #{GoVersion}
	
	div
		h2 pigen.dk
		form#server-sickbeard-restart
			[method="POST"]
		input.btn[type="submit"]
			[value="Genstart SickBeard"]

	script[type="text/javascript"]
		#{Script}
`

const ServerTemplateScript = `
		var checks = {
				dropbox: ["/dropbox/status", "dropbox status", 2000, "span6"],
				landscapeSysinfo: ["/landscape/sysinfo", "landscape sysinfo", 2000, "span8"],
				dstat: ["/dstat", "dstat 1 10", 20000, "span6"],
				dropboxHelp: ["/dropbox/help", "dropbox help", 2000, "span6"],
		};

		$(function () {

			var $form = $("form#server-sickbeard-restart"),
				$button = $form.find("input:submit");

			// Handle form submit
			$form.submit(function (e) {

				$button.removeClass("btn-inverse");
				$button.addClass("disabled");
				$button.attr("disabled", "disabled");

				$.post("/action", {action: "server-sickbeard-restart"}).always(function () {
					$button.removeClass("disabled");
					$button.removeAttr("disabled");
					$button.addClass("btn-inverse");
				});
				e.preventDefault()
				return false;
			});

				var $body = $(document.body),
					pre, data;

				for (var id in checks) {
					data = checks[id];
					pre = $('<div id="'+id+'"><h2>'+data[1]+'</h2><pre>loading...</pre></div>').appendTo(document.body).find("pre");

					(function  (p, id, data) {
						var f = function() {
							p.load(data[0]);
						};
						setInterval(f, data[2]);
					})(pre, id, data);
				}
			});

			// setTimeout(function () {window.location.reload();}, 20000);
`
