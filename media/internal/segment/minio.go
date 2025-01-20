package segment

import (
	"context"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/odit-bit/sone/media/internal/domain"
	"github.com/odit-bit/sone/pkg/buffer"
)

var _ domain.SegmentStorer = (*minioStore)(nil)
var _ domain.SegmentFetcher = (*minioStore)(nil)

type minioStore struct {
	mio    *minio.Client
	bucket string
}

func NewMiniStore(ctx context.Context, bucket string, mio *minio.Client) (*minioStore, error) {
	if ok, _ := mio.BucketExists(ctx, bucket); !ok {
		err := mio.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}
	return &minioStore{
		mio:    mio,
		bucket: bucket,
	}, nil
}

// GetVideoSegment implements domain.LiveRepository.
func (m *minioStore) GetVideoSegment(ctx context.Context, path string) (*domain.Segment, error) {
	// path := path.Join(videoID, part)
	obj, err := m.mio.GetObject(ctx, m.bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return &domain.Segment{Body: obj}, nil
}

func (m *minioStore) InsertVideoSegment(_ context.Context, path string, content io.Reader) error {
	// path := path.Join(videoID, part)

	b := buffer.Pool.Get()
	n, err := io.Copy(b, content)
	if err != nil {
		return nil
	}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()
	_, err = m.mio.PutObject(ctx, m.bucket, path, b, n, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	buffer.Pool.Put(b)
	return nil
}
