<html>
	<head>
		<title>Network Analysis - Path</title>
	</head>
	<body>
		<h1>Network Analysis - Path</h1>
		<p>Full path: {{.Fullpath}}</p>
		<p>Cost: {{.Cost}}</p>
		<div id="replpath"></div>
		<script type="text/javascript">
xhrpath = new XMLHttpRequest();
xhrpath.onreadystatechange = function() { if(xhrpath.readyState == 4) {
	if(xhrpath.status == 200) {
		document.querySelector("#replpath").innerHTML = xhrpath.responseText;
	} else {
		document.querySelector("#replpath").innerHTML = '<p style="color:red">' + xhrpath.statusText + "</p>";
	}
}};
xhrpath.open("GET", "/render/path/{{.Dbname}}/{{.Netpath}}", true);
xhrpath.send();
		</script>
	</body>
</html>
