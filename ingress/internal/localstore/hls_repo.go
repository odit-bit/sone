package localstore

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/odit-bit/sone/ingress/internal/av"
	"github.com/odit-bit/sone/ingress/internal/domain"
	"github.com/spf13/afero"
)

var _ domain.HLSRepository = (*hlsRepository)(nil)

type hlsRepository struct {
	root    afero.Fs
	rootDir string
	skRepo  domain.StreamKeyRepository //*kvstore.Client
}

func NewHLSRepository(root string, afs afero.Fs, skRepo domain.StreamKeyRepository) *hlsRepository {
	// afs := afero.NewBasePathFs(afero.NewOsFs(), root)
	return &hlsRepository{
		root:    afs,
		rootDir: root,
		skRepo:  skRepo,
	}
}

// TranscodeFLV implements domain.HLSRepository.
func (h *hlsRepository) TranscodeFLV(ctx context.Context, stream_key, stream_name string, r io.Reader) <-chan *domain.OpError {
	respC := make(chan *domain.OpError, 1)

	if !h.skRepo.Exist(stream_key) {
		respC <- domain.NewOpError("stream key not exist")
		close(respC)
		return respC
	}

	writeErr := func(err error) {
		opErr := domain.NewOpError(err.Error())
		respC <- opErr
		close(respC)
	}

	go func() {
		// path for hls file
		// root/stream_key/stream_name/index.m3u8..segmentN.ts
		dir := filepath.Join(stream_key, stream_name)
		if ok, err := afero.Exists(h.root, dir); err != nil {
			writeErr(err)
			return
		} else if ok {
			err = errors.Join(err, fmt.Errorf("already streaming"))
			writeErr(err)
			return
		}

		err := h.root.MkdirAll(dir, 0775)
		if err != nil {
			writeErr(err)
			return
		}

		streamPath := filepath.Join(h.rootDir, dir)
		err = av.WriteFLVToHLS(ctx, filepath.Join(streamPath, "index.m3u8"), r)
		if err != nil {
			err = errors.Join(err, h.deleteContents(stream_key))
			writeErr(err)
			return
		}
	}()

	return respC

}

func (h *hlsRepository) deleteContents(dirPath string) error {
	info, err := afero.ReadDir(h.root, dirPath)
	if err != nil {
		return err
	}
	for _, i := range info {
		path := filepath.Join(dirPath, i.Name())
		if err := h.root.RemoveAll(path); err != nil {
			return err
		}
	}

	return err
}
