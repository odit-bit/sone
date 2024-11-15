package fs

import (
	"context"
	"io/fs"
	"sync"

	"github.com/odit-bit/sone/internal/kvstore"
	"github.com/odit-bit/sone/stream/internal/domain"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

var _ domain.StreamKeyRepo = (*FS)(nil)
var _ domain.EntryRepository = (*FS)(nil)

// handle persistence of hls stream media
type FS struct {
	mx sync.Mutex

	root afero.Fs

	keyRepo *kvstore.Client
	logger  *logrus.Logger
}

func New(afs afero.Fs, kv *kvstore.Client) (*FS, error) {
	return &FS{
		mx:      sync.Mutex{},
		root:    afs,
		keyRepo: kv,
		logger:  logrus.StandardLogger(),
	}, nil
}

func (h *FS) GetSegment(name string, segment *domain.Segment) error {
	f, err := h.root.Open(name)
	if err != nil {
		return err
	}
	segment.Body = f
	return nil
}

// Exist implements domain.StreamKeyRepo.
func (h *FS) Exist(ctx context.Context, key domain.Key) bool {
	ok := h.keyRepo.Exist(string(key))
	return ok
}

// Save implements domain.StreamKeyRepo.
func (h *FS) Save(ctx context.Context, key domain.Key) error {
	h.keyRepo.Set(string(key), "")
	return nil
}

// func (h *FS) GetSegmentOld(w io.Writer, file string) error {

// 	f, err := h.root.Open(file)
// 	if err != nil {
// 		// http.Error(w, err.Error(), http.StatusNotFound)
// 		return err
// 	}
// 	defer f.Close()

// 	if _, err := io.Copy(w, f); err != nil {
// 		return err
// 	}
// 	return nil
// }

func (h *FS) List(offset, limit int) <-chan domain.Entry {
	eC := make(chan domain.Entry)

	count := 0
	if limit < 1 {
		limit = 1000
	}
	go func() {
		err := afero.Walk(h.root, ".", func(path string, d fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if path == "." {
				return nil
			}

			if count < limit {
				if d.IsDir() {
					eC <- domain.Entry{
						Name: d.Name(),
					}
					count++
				}
				return fs.SkipDir
			}
			return fs.SkipAll
		})
		if err != nil {
			h.logger.Error(err)
		}

		close(eC)
	}()

	return eC
}

////////////////////////////////////////////////////

func (h *FS) SetLogger(l *logrus.Logger) {
	h.mx.Lock()
	h.logger = l
	h.mx.Unlock()
}

// func (h *FileHandler) View(w http.ResponseWriter, r *http.Request) {
// 	source := r.URL.Query().Get("source")
// 	if source == "" {
// 		http.Error(w, "empty source file", http.StatusBadRequest)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "text/html")
// 	w.WriteHeader(http.StatusOK)

// 	index := filepath.Join("file", source, "index.m3u8")
// 	writeVideoPlayer(w, index)

// }

// func (h *FileHandler) File(w http.ResponseWriter, r *http.Request) {

// 	dir := r.PathValue("dir")
// 	segment := r.PathValue("segment")
// 	filename := filepath.Join(dir, segment)
// 	log.Println("stream path:", dir, segment)
// 	if filename == "" {
// 		http.NotFound(w, r)
// 		return
// 	}
// 	f, err := h.root.Open(filename)
// 	if err != nil {
// 		http.NotFound(w, r)
// 		return
// 	}
// 	defer f.Close()

// 	stat, err := f.Stat()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	ff, ok := f.(io.ReadSeeker)
// 	if !ok {
// 		err := fmt.Errorf("file [%s] not implement io.ReadSeeker ", filename)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	http.ServeContent(w, r, stat.Name(), stat.ModTime(), ff)
// }
