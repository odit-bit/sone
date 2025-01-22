###  
Sone is video streaming server

###  Getting Started
Before to use , make sure the environment (os) has ffmpeg installed.see [install ffmpeg](https://ffmpeg.org/download.html)


#### Build from source
```
git clone https://github.com/odit-bit/sone.git
cd sone
go build -o ./bin/sone ./cmd/sone
```

run server
```
./bin/sone server run --debug
```


#### streaming (RTMP)


##### OBS
Set the destination server to `rtmp://localhost:1935`,set stream-key that generate from server .
start the streaming and visit [http://localhost:9797](http://localhost:9797/stream) to playback the stream 

##### ffmpeg
```
ffmpeg -re -i input.mp4 -c:v libx264 -preset veryfast -b:v 2800k -maxrate 2996k -bufsize 4200k -pix_fmt yuv420p -g 50 -c:a aac -b:a 128k -ac 2 -ar 44100 -f flv rtmp://localhost:1935/live/{stream-key}
```


<!-- #### streaming (hls)

```
 ffmpeg -re -i test-stream.mp4 -codec copy -f hls  -hls_list_size 0 http://localhost:6969/{stream-key}/{playlist}.m3u8
``` -->

### HOW
sone has 2 module ingest and playback.  
Server (rtmp) will receive rtmp protocol, serialized (to flv), transcode (ffmpeg) to playlist segment `.ts` file and store at disk (cache) before available for HLS distribution on http, it will store up to last 6 second of the stream for playback.



