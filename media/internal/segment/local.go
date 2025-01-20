package segment

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/odit-bit/sone/media/internal/domain"
	"github.com/spf13/afero"
)

var _ domain.SegmentRepository = (*localStore)(nil)

type localStore struct {
	blob afero.Fs
}

func NewLocalFS(fs afero.Fs) *localStore {
	return &localStore{
		blob: fs,
	}
}

// GetVideoSegment implements domain.SegmentRepository.
func (l *localStore) GetVideoSegment(ctx context.Context, key string) (*domain.Segment, error) {
	// log.Printf("\n get: %v \n", key)
	f, err := l.blob.Open(key)
	if err != nil {
		return nil, err
	}
	return &domain.Segment{Body: f}, nil
}

// InsertVideoSegment implements domain.SegmentRepository.
func (l *localStore) InsertVideoSegment(ctx context.Context, key string, r io.Reader) error {

	// log.Printf("\n put: %v \n", key)
	f, err := l.blob.OpenFile(key, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}
	log.Println("write media:", n, "path:", key)
	return f.Close()

}
