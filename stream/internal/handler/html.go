package handler

type VideoPlayerData struct {
	SourceURL string
}

const VideoPlayerHTML = `
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

                hls.loadSource({{.SourceURL}});
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

var VideoListHTML = `
<!DOCTYPE html>
<html>
<body>

<h2>Current Streaming</h2>


<ol>
{{range .}} 
  <li> <a href={{.Path}}>{{.Name}}</a> </li>
{{end}}
</ol>  


</body>
</html>
`

// var VideoListTemplate = func() *template.Template {
// 	t, err := template.New("listPage").Parse(VideoListHTML)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return t
// }()

// type ListPageData struct {
// }

// func writeListPage(w io.Writer, entry []Entry) {
// 	VideoListTemplate.Execute(w, entry)
// }
