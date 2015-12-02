<html>
	<head>
		<title>Network Analysis - Databases</title>
	</head>
	<body>
		<h1>Network Analysis - Databases</h1>
		<p>The following databases are available on this server:</p>
		<ul>
			{{range .}}<li><a href="/db/{{.}}">{{.}}</a></li>{{end}}
		</ul>
	</body>
</html>
