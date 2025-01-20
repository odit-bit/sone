package web

import (
	"io"
	"text/template"
)

type CreateStreamArgs struct {
	ApiEndpoint string
}

var CreateStreamPage = `
<!DOCTYPE html>
<html>
<body>

<h1>Create Stream</h1>

<form action={{.ApiEndpoint}} method="post">
  <label for="fname">stream name:</label>
  <input type="text" id="fname" name="fname"><br><br>
  <input type="submit" value="Submit">
</form>

</body>
</html>

`

type CreateStreamTemplate struct {
	data CreateStreamArgs
	tmpl *template.Template
}

func NewCreateStreamTemplate(apiEndpoint string) *CreateStreamTemplate {
	t := template.Must(template.New("create_stream_page").Parse(CreateStreamPage))
	return &CreateStreamTemplate{
		data: CreateStreamArgs{
			ApiEndpoint: apiEndpoint,
		},
		tmpl: t,
	}
}

func (t *CreateStreamTemplate) Render(w io.Writer, data any) error {
	if data == nil {
		data = t.data
	}
	return t.tmpl.Execute(w, data)
}

////////////////////////////
//  Stream Root

type StreamRootArgs struct {
	List           []streamListArgs
	CreateStreamCB string
}
type streamListArgs struct {
	PlaybackUrl string
	Title       string
}

var StreamRootHTML = `
<!DOCTYPE html>
<html>
<body>

	<form action={{.CreateStreamCB}} method="post">
		<!-- <label for="fname">stream name:</label> -->
		<!-- <input type="text" id="fname" name="fname"><br><br> -->
		<input type="submit" value="stream-key">
	</form>

	<h2>Current Streaming</h2>

	<ol>
		{{range .List}} 
			<li> <a href={{.PlaybackUrl}}>{{.Title}}</a> </li>
		{{end}}
	</ol>  


</body>
</html>
`

type StreamPageTemplate struct {
	data StreamRootArgs
	tmpl *template.Template
}

func NewStreamRootTemplate() *StreamPageTemplate {
	tmpl := template.Must(template.New("streamRoot").Parse(StreamRootHTML))
	return &StreamPageTemplate{
		data: StreamRootArgs{},
		tmpl: tmpl,
	}
}

func (t *StreamPageTemplate) Render(w io.Writer, data any) error {
	if data == nil {
		data = t.data
	}
	return t.tmpl.Execute(w, data)
}

/////////////////////////// Stream Playback

type PlaybackData struct {
	SourceURL string
}

var PlaybackHTML = `
<!DOCTYPE html>
<html>

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width">
    <style>
        html,
        body {
            margin: 0;
            padding: 0;
            height: 100%;
            overflow: hidden;
        }

        #video {
            width: 100%;
            height: 100%;
            background: black;
        }
    </style>
</head>

<body>

    <video id="video" muted controls autoplay playsinline></video>

    <script src="https://cdn.jsdelivr.net/npm/hls.js@1.2.9"></script>

    <script>

        const create = () => {
            const video = document.getElementById('video');

            // always prefer hls.js over native HLS.
            // this is because some Android versions support native HLS
            // but don't support fMP4s.
            if (Hls.isSupported()) {
                const hls = new Hls({
                    maxLiveSyncPlaybackRate: 1.5,
                });

                hls.on(Hls.Events.ERROR, (evt, data) => {
                    if (data.fatal) {
                        hls.destroy();

                        setTimeout(create, 2000);
                    }
                });

				//make sure the source with the (" ")
                hls.loadSource("{{.SourceURL}}");
                hls.attachMedia(video);

                video.play();

            } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
                // since it's not possible to detect timeout errors in iOS,
                // wait for the playlist to be available before starting the stream
                fetch('index.m3u8')
                    .then(() => {
                        video.src = 'index.m3u8';
                        video.play();
                    });
            }
        };

        window.addEventListener('DOMContentLoaded', create);

    </script>

</body>

</html>
`

type PlaybackTemplate struct {
	data PlaybackData
	tmpl *template.Template
}

func NewPlaybackTemplate() *PlaybackTemplate {
	tmpl := template.Must(template.New("StreamPlayback").Parse(PlaybackHTML))
	return &PlaybackTemplate{
		data: PlaybackData{},
		tmpl: tmpl,
	}
}

func (t *PlaybackTemplate) Render(w io.Writer, data any) error {
	if data == nil {
		data = t.data
	}
	return t.tmpl.Execute(w, data)
}
