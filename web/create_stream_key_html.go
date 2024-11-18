package web

// var VideoPlayerTemplate = func() *template.Template {
// 	t, err := template.New("VideoPlayer").Parse(VideoPlayerHTML)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return t
// }()

// func writeVideoPlayer(w io.Writer, playlistURL string) {
// 	ipd := &VideoPlayerData{
// 		SourceURL: playlistURL,
// 	}
// 	VideoPlayerTemplate.Execute(w, ipd)
// }


var StreamKeyCreateHTML = `
<!DOCTYPE html>
<html>
<body>

<h1>The input element</h1>

<form action={{.CreateStreamKeyURL}}>
  <label for="fname">stream name:</label>
  <input type="text" id="fname" name="fname"><br><br>
  <input type="submit" value="Submit">
</form>

</body>
</html>

`
