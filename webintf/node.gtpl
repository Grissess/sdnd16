<html>
	<head>
		<title>Network Analysis - Node</title>
	</head>
	<body>
		<h1>Network Analysis - Node</h1>
		<p>Database: {{.Dbname}}</p>
		<p>Node: {{.Node}}</p>
		<div id="replnode"></div>
		<script type="text/javascript">
xhrnode = new XMLHttpRequest();
xhrnode.onreadystatechange = function() { if(xhrnode.readyState == 4) {
	if(xhrnode.status == 200) {
		document.querySelector("#replnode").innerHTML = xhrnode.responseText;
	} else {
		document.querySelector("#replnode").innerHTML = '<p style="color:red">' + xhrnode.statusText + "</p>";
	}
}};
xhrnode.open("GET", "/render/node/{{.Dbname}}/{{.Node}}", true);
xhrnode.send();
		</script>
	</body>
</html>
