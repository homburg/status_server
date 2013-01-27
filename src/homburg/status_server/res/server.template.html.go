package status_server

const ServerTemplate = `<!DOCTYPE HTML>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Status page</title>
	<style type="text/css">
		div {
			margin: 1em;
		}
	</style>
	<script type="text/javascript" src="http://code.jquery.com/jquery-latest.js"></script>
	<link href="//netdna.bootstrapcdn.com/twitter-bootstrap/2.2.2/css/bootstrap-combined.min.css" rel="stylesheet">
</head>
<body>
	<div><h1>{{.}}</h1></div>

	<script type="text/javascript">
		var checks = {
				dropbox: ["/dropbox/status", "dropbox status", 2000, "span6"],
				landscapeSysinfo: ["/landscape/sysinfo", "landscape sysinfo", 2000, "span8"],
				dstat: ["/dstat", "dstat 1 10", 20000, "span6"],
				dropboxHelp: ["/dropbox/help", "dropbox help", 2000, "span6"],
		};

		$(function () {
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
	</script>
</body>
</html>`
