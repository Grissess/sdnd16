<html>
	<head>
		<title>Network Analysis - Database</title>
	</head>
	<body>
		<h1>Network Analysis - Database</h1>
		<p>Database: {{.Dbname}}</p>
		<div id="repldb"></div>
		<script type="text/javascript">
xhrdb = new XMLHttpRequest();
xhrdb.onreadystatechange = function() { if(xhrdb.readyState == 4) {
	if(xhrdb.status == 200) {
		document.querySelector("#repldb").innerHTML = xhrdb.responseText;
	} else {
		document.querySelector("#repldb").innerHTML = '<p style="color:red">' + xhrdb.statusText + "</p>";
	}
}};
xhrdb.open("GET", "/render/db/{{.Dbname}}", true);
xhrdb.send();
		</script>
	</body>
</html>
