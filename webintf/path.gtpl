<html>
	<head>
		<title>Network Analysis - Path</title>
	</head>
	<body>
		<h1>Network Analysis - Path</h1>
		<p>Raw result: {{.Rawpath}}</p>
		<p>Full path: {{range .Path}}{{.}} -> {{end}}</p>
		<p>Cost: {{.Cost}}</p>
		<img src="/render/path/{{.Dbname}}/{{.Netpath}}" alt="SVG path"/>
	</body>
</html>
