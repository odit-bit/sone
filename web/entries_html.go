package web

type Entry struct {
	Endpoint string
	Name     string
}

var VideoListHTML = `
<!DOCTYPE html>
<html>
<body>

<h2>Current Streaming</h2>


<ol>
{{range .}} 
  <li> <a href={{.Endpoint}}>{{.Name}}</a> </li>
{{end}}
</ol>  


</body>
</html>
`
