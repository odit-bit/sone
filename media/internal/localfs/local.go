package localfs

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/odit-bit/sone/media/internal/domain"
	"github.com/odit-bit/sone/pkg/kvstore"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// var _ domain.LiveRepository = (*local)(nil)
// var _ domain.EntryRepository = (*local)(nil)
// var _ domain.StreamKeyRepository = (*local)(nil)

type local struct {
	root   afero.Fs
	kv     *kvstore.Client
	logger *logrus.Logger
}

func New(afs afero.Fs, kv *kvstore.Client, logger *logrus.Logger) *local {
	return &local{
		root:   afs,
		kv:     kv,
		logger: logger,
	}
}

// // Exist implements domain.StreamKeyRepository.
// func (h *local) Exist(ctx context.Context, key string) bool {
// 	return h.kv.Exist(string(key))
// }

// // Save implements domain.StreamKeyRepository.
// func (h *local) Save(ctx context.Context, key string) error {
// 	return h.kv.Set(string(key), "")
// }

// Get implements domain.SegmentRepository.
func (h *local) GetVideoSegment(ctx context.Context, key, sequence string) (*domain.Segment, error) {
	segmentPath := filepath.Join(key, sequence)
	var segment domain.Segment
	f, err := h.root.Open(segmentPath)
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
		defer close(eC)
		h.logger.Debug(fmt.Sprintf("start iterate stream list start at %v, limit %v ", offset, limit))
		rootDir := "."
		f, err := h.root.Open(rootDir)
		if err != nil {
			h.logger.Error(err)
			return
		}
		defer f.Close()

		entries, err := f.Readdir(limit)
		if err != nil {
			h.logger.Error(err)
			return
		}

		for _, e := range entries {
			if e.IsDir() {
				h.logger.Debug(fmt.Sprintf("found entry, %v", e.Name()))
				eC <- domain.Entry{
					Endpoint: e.Name(),
					Name:     e.Name(),
				}
				// subPath := filepath.Join(rootDir, e.Name())
				// f2, err := h.root.Open(subPath)
				// if err != nil {
				// 	//end of dir
				// 	f2.Close()
				// 	continue
				// }
				// infos, err := f2.Readdir(limit)
				// if err != nil {
				// 	// h.logger.Error(err)
				// 	f2.Close()
				// 	continue
				// }
				// for _, info := range infos {
				// 	if info.IsDir() {
				// 		eC <- domain.Entry{
				// 			Endpoint: filepath.Join(subPath, info.Name()),
				// 			Name:     info.Name(),
				// 		}
				// 	}
				// }
				// f2.Close()
			}
		}
		h.logger.Debug("iterate stream list finish")

	}()
	return eC
}
