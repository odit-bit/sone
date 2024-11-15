package av

import (
	"fmt"
	"io"
	"log"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func PresetLightNormalizePreview(in io.Reader, out io.Writer) <-chan error {
	log.Println("startProcessing preview")
	done := make(chan error)
	go func() {
		// input pipe
		err := ffmpeg.Input("pipe:",
			ffmpeg.KwArgs{"t": "10"}).Audio().
			// ffmpeg.KwArgs{}).Audio().
			Filter("crystalizer", ffmpeg.Args{}, ffmpeg.KwArgs{}).
			Filter("loudnorm", ffmpeg.Args{}, ffmpeg.KwArgs{"tp": "-0.1", "lra": 5, "i": "-16.0"}).
			WithInput(in).
			Output("pipe:", ffmpeg.KwArgs{"f": "mp3"}).
			WithOutput(out).
			ErrorToStdOut().
			Run()
		done <- err
		close(done)
	}()
	return done
}

func PresetLightNormalize(in io.Reader, out io.Writer) <-chan error {
	log.Println("startProcessing preview")
	done := make(chan error)
	go func() {
		// input pipe
		err := ffmpeg.Input("pipe:",
			ffmpeg.KwArgs{}).Audio().
			Filter("crystalizer", ffmpeg.Args{}, ffmpeg.KwArgs{}).
			Filter("loudnorm", ffmpeg.Args{}, ffmpeg.KwArgs{"tp": "-0.1", "lra": 5, "i": "-16.0"}).
			WithInput(in).
			Output("pipe:", ffmpeg.KwArgs{"f": "mp3"}).
			WithOutput(out).
			ErrorToStdOut().
			Run()
		done <- err
		close(done)
	}()
	return done
}

/* the next couple function is preceded preset that breakdown into separate function */

type ProcessorFunc func(*ffmpeg.Stream) *ffmpeg.Stream
type OutFormatterFunc func(*ffmpeg.Stream, io.Writer) *ffmpeg.Stream

// process full audio content
func HandleFull(in io.Reader, out io.Writer, formatterFn OutFormatterFunc, processors ...ProcessorFunc) error {
	// input pipe
	audio := ffmpeg.Input("pipe:", ffmpeg.KwArgs{})
	for _, p := range processors {
		audio = p(audio)
	}
	// // input
	audio = audio.WithInput(in)
	finish := formatterFn(audio, out)

	if err := finish.ErrorToStdOut().Run(); err != nil {
		return fmt.Errorf("sone error process stream: %v", err)
	}

	return nil
}

// process only 10 second of audio content
func HandlePreview(in io.Reader, out io.Writer, formatterFn OutFormatterFunc, processors ...ProcessorFunc) error {
	// input pipe
	audio := ffmpeg.Input("pipe:", ffmpeg.KwArgs{"t": "10"})
	for _, p := range processors {
		audio = p(audio)
	}
	// // input
	audio = audio.WithInput(in)
	finish := formatterFn(audio, out)

	if err := finish.ErrorToStdOut().Run(); err != nil {
		return fmt.Errorf("sone error process stream: %v", err)
	}
	return nil
}

func LightNormalizeStreamProcessor() func(stream *ffmpeg.Stream) *ffmpeg.Stream {
	return func(stream *ffmpeg.Stream) *ffmpeg.Stream {
		filter := stream.
			//noise sharp
			Filter("crystalizer", ffmpeg.Args{}, ffmpeg.KwArgs{}).
			//norm
			Filter("loudnorm", ffmpeg.Args{}, ffmpeg.KwArgs{"tp": "-0.1", "lra": 5, "i": "-16.0"})

		return filter
	}
}

func WavStreamOut(stream *ffmpeg.Stream, out io.Writer) *ffmpeg.Stream {
	// output
	output := stream.Output("pipe:", ffmpeg.KwArgs{"f": "wav"}).WithOutput(out)
	return output
}

func MP3StreamOut(stream *ffmpeg.Stream, out io.Writer) *ffmpeg.Stream {
	// output
	output := stream.Output("pipe:", ffmpeg.KwArgs{"f": "mp3"}).WithOutput(out)
	return output
}
