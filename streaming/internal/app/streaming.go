package app

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/odit-bit/sone/streaming/internal/api"
	"github.com/odit-bit/sone/streaming/internal/database"
	"github.com/odit-bit/sone/streaming/internal/scheme"
	"github.com/segmentio/ksuid"
)

type Application struct {
	rpc *api.GrpcClient
	db  *database.Sqlite
}

func New(db *database.Sqlite, usersRepo *api.GrpcClient) *Application {
	return &Application{
		rpc: usersRepo,
		db:  db,
	}
}

type InsertStreamOpt struct {
	Title string
}

func (a *Application) InsertStream(ctx context.Context, token string, opt InsertStreamOpt) (*scheme.LiveStream, error) {
	if err := a.validateInsertStream(&opt); err != nil {
		return nil, err
	}

	//auth user
	_, err := a.rpc.Users.AuthToken(ctx, token)
	if err != nil {
		return nil, err
	}

	sk, err := generateRandomString(12)
	if err != nil {
		return nil, err
	}

	ls := scheme.LiveStream{
		ID:     generateID(),
		Title:  opt.Title,
		IsLive: 0,

		Key: sk,
	}

	tx, err := a.db.LiveStream.Save(ctx, &ls)
	if err != nil {
		return nil, err
	}

	return &ls, tx.Commit()
}

func (a *Application) StartStream(ctx context.Context, key string) (*scheme.LiveStream, error) {
	c, err := a.db.LiveStream.FindFilter(ctx, scheme.Filter{Key: key})
	if err != nil {
		return nil, fmt.Errorf("streaming: %v", err)
	}

	ls, ok := c.Next()
	if !ok {
		return nil, err
	}

	ls.IsLive = 1
	tx, err := a.db.LiveStream.Save(ctx, ls)
	if err != nil {
		return nil, err
	}

	return ls, tx.Commit()
}

func (a *Application) EndStream(ctx context.Context, key string) (*scheme.LiveStream, error) {
	c, err := a.db.LiveStream.FindFilter(ctx, scheme.Filter{Key: key})
	if err != nil {
		return nil, fmt.Errorf("streaming: %v", err)
	}

	ls, ok := c.Next()
	if !ok {
		return nil, err
	}
	// set live false
	ls.IsLive = 0
	// save

	tx, err := a.db.LiveStream.Save(ctx, ls)
	if err != nil {
		return nil, err
	}

	return ls, tx.Commit()
}

func (a *Application) ListStream(ctx context.Context, id string) (*scheme.LiveStreamsErr, error) {
	var filter scheme.Filter
	if id == "" {
		filter.Id = id
	}
	c, err := a.db.LiveStream.FindFilter(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("streaming: record not found")
	}
	return c, nil
}

// HELPER

func (a *Application) validateInsertStream(opt *InsertStreamOpt) error {
	if opt.Title == "" {
		title, _ := generateRandomString(6)
		opt.Title = fmt.Sprintf("video%s", title)
	}

	return nil
}

func generateID() string {
	return ksuid.New().String()
}

// func GenerateKey() string {
// 	return generateKey()
// }

func generateKey() string {
	var b string
	count := 0
	for {
		v, err := generateRandomString(12)
		if err != nil {
			if count < 3 {
				count++
				continue
			} else {
				panic("failed generate random string more than 3 times")
			}
		}
		b = v
		break
	}

	return b
}

func generateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[randomIndex.Int64()]
	}
	return string(b), nil
}
