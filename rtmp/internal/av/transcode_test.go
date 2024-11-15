package av

import (
	"context"
	"os"
	"testing"

	"github.com/spf13/afero"
)

func TestXxx(t *testing.T) {
	root := "/mnt/d/wsl/test-stream"
	afs := afero.NewBasePathFs(afero.NewOsFs(), root)
	ok, err := afero.DirExists(afs, root)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Skip("directory not exist ")
	}

	if err := afs.MkdirAll("test", 0666); err != nil {
		t.Fatal(err)
	}

	errC := make(chan error, 1)
	f, err := os.OpenFile("/mnt/d/wsl/test.flv", os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	go func() {
		if err := WriteFLVToHLS(context.TODO(), "/mnt/d/wsl/test-stream/index.m3u8", f); err != nil {
			errC <- err
		}
		close(errC)
	}()

	err, ok = <-errC
	if ok {
		t.Fatal(err)
	}
}
