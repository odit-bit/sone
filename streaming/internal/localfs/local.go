package localfs

import (
	"context"
	"path/filepath"

	"github.com/odit-bit/sone/pkg/kvstore"
	"github.com/odit-bit/sone/streaming/internal/domain"
	"github.com/spf13/afero"
)

var _ domain.SegmentRepository = (*local)(nil)
var _ domain.EntryRepository = (*local)(nil)
var _ domain.StreamKeyRepository = (*local)(nil)

type local struct {
	root afero.Fs
	kv   *kvstore.Client
}

// Exist implements domain.StreamKeyRepository.
func (h *local) Exist(ctx context.Context, key domain.Key) bool {
	return h.kv.Exist(string(key))
}

// Save implements domain.StreamKeyRepository.
func (h *local) Save(ctx context.Context, key domain.Key) error {
	return h.kv.Set(string(key), "")
}

func New(afs afero.Fs, kv *kvstore.Client) *local {
	return &local{
		root: afs,
		kv:   kv,
	}
}

// Get implements domain.SegmentRepository.
func (h *local) Get(ctx context.Context, name string) (*domain.Segment, error) {
	var segment domain.Segment
	f, err := h.root.Open(name)
	if err != nil {
		return nil, err
	}
	segment.Body = f
	return &segment, nil
}

// List implements domain.EntryRepository.
func (h *local) List(ctx context.Context, offset int, limit int) <-chan domain.Entry {
	eC := make(chan domain.Entry)

	// count := 0
	if limit < 1 {
		limit = 1000
	}

	go func() {
		rootDir := "."
		f, err := h.root.Open(rootDir)
		if err != nil {
			// h.logger.Error(err)
			return
		}
		defer f.Close()

		entries, err := f.Readdir(limit)
		if err != nil {
			// h.logger.Error(err)
			return
		}

		for _, e := range entries {
			if e.IsDir() {
				subPath := filepath.Join(rootDir, e.Name())
				f2, err := h.root.Open(subPath)
				if err != nil {
					//end of dir
					f2.Close()
					continue
				}
				infos, err := f2.Readdir(limit)
				if err != nil {
					// h.logger.Error(err)
					f2.Close()
					continue
				}
				for _, info := range infos {
					if info.IsDir() {
						eC <- domain.Entry{
							Endpoint: filepath.Join(subPath, info.Name()),
							Name:     info.Name(),
						}
					}
				}
				f2.Close()
			}
		}

		close(eC)
	}()
	return eC
}
