package application

import (
	"context"

	"github.com/odit-bit/sone/users/internal/domain"
)

type Application struct {
	userRepo   domain.UserRepository
	streams    domain.StreamRepository
	userGlRepo domain.UserGoogleRepository
}

func New(glusers domain.UserGoogleRepository, users domain.UserRepository, streams domain.StreamRepository) *Application {
	return &Application{
		userRepo:   users,
		streams:    streams,
		userGlRepo: glusers,
	}
}

type RegisterUserArgs struct {
	Username string
	Password string
}

type RegisterUserResponse struct {
	Id    string
	Token string
}

func (a *Application) RegisterUser(ctx context.Context, args RegisterUserArgs) (RegisterUserResponse, error) {
	var resp RegisterUserResponse
	user, err := domain.CreateUser(args.Username, args.Password)
	if err != nil {
		return resp, err
	}

	if err := a.userRepo.Save(ctx, user); err != nil {
		return resp, err
	}
	resp.Id = user.ID
	resp.Token = user.Token
	return resp, nil
}

type AuthenticateUserArgs struct {
	Username string
	Password string
}

func (a *Application) AuthenticateUser(ctx context.Context, args AuthenticateUserArgs) (string, error) {
	user, err := a.userRepo.Find(ctx, domain.FilterOption{Field: domain.Name_Field, Value: args.Username})
	if err != nil {
		return "", err
	}
	if err := user.Authenticate(args.Password); err != nil {
		return "", err
	}

	token := domain.GenerateKey()
	user.Token = token
	if err := a.userRepo.Save(ctx, &user); err != nil {
		return "", err
	}
	return token, nil
}

type UserInfo struct {
	ID       string
	Username string
	Email    string
}

func (a *Application) AuthToken(ctx context.Context, token string) (UserInfo, error) {
	user, err := a.userRepo.Find(ctx, domain.FilterOption{Field: "token", Value: token})
	if err != nil {
		return UserInfo{}, err
	}

	return UserInfo{ID: user.ID, Username: user.Name}, nil
}

// type CreateStreamArgs struct {
// 	Token string
// 	Title string
// }

// type CreateStreamResponse struct {
// 	Key         string
// 	ID          string
// 	Title       string
// 	PublishedAt time.Time
// 	IsLive      bool
// }

// func (a *Application) InsertStream(ctx context.Context, args CreateStreamArgs) (domain.Stream, error) {
// 	if _, err := a.userRepo.Find(ctx, domain.FilterOption{Field: domain.Token_Field, Value: args.Token}); err != nil {
// 		log.Println("error", err)
// 		return domain.Stream{}, fmt.Errorf("user not found")

// 		// if _, ok, _ := a.userGlRepo.Get(ctx, args.UserID); !ok {
// 		// 	return CreateStreamResponse{}, fmt.Errorf("user not found")
// 		// }
// 	}

// 	stream := domain.NewStream(args.Title)

// 	if err := a.streams.Save(ctx, stream); err != nil {
// 		return domain.Stream{}, err
// 	}
// 	return stream, nil
// }

// type StartStreamArgs struct {
// 	StreamKey string
// }

// func (a *Application) StartStream(ctx context.Context, args StartStreamArgs) error {
// 	stream, ok, err := a.streams.Get(ctx, domain.StreamGetOption{Key: args.StreamKey})
// 	if !ok {
// 		if err != nil {
// 			// logging the error
// 		}
// 		return fmt.Errorf("invalid stream key")
// 	}
// 	if stream.IsStreaming {
// 		// logging it , streaming key should use only one stream
// 		return fmt.Errorf("already streaming")
// 	} else {
// 		stream.IsStreaming = true
// 		stream.LastStartStream = time.Now()
// 	}
// 	if err := a.streams.Save(ctx, stream); err != nil {
// 		return err
// 	}
// 	return nil
// }

// type EndStreamArgs struct {
// 	StreamKey string
// }

// func (a *Application) EndStream(ctx context.Context, args EndStreamArgs) error {
// 	stream, ok, err := a.streams.Get(ctx, domain.StreamGetOption{Key: args.StreamKey})
// 	if !ok {
// 		if err != nil {
// 			// logging the error
// 		}
// 		return fmt.Errorf("invalid stream key")
// 	}
// 	if !stream.IsStreaming {
// 		// logging it , this is a bug
// 		return fmt.Errorf("is not streaming")
// 	}
// 	stream.IsStreaming = false
// 	stream.LastEndStream = time.Now()
// 	if err := a.streams.Save(ctx, stream); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (a *Application) GetStream(ctx context.Context, id string) (domain.Stream, bool, error) {
// 	s, ok, err := a.streams.Get(ctx, domain.StreamGetOption{ID: id})
// 	if !ok {
// 		return s, false, err
// 	}
// 	return s, true, nil
// }

// func (a *Application) ListStream(ctx context.Context, limit, offset int) (<-chan domain.Stream, error) {
// 	return a.streams.List(ctx, limit), nil
// }

// google user

func (a *Application) GetUserGoogle(ctx context.Context, id string) (domain.UserGl, bool, error) {
	return a.userGlRepo.Get(ctx, id)
}

func (a *Application) SaveUserGoogle(ctx context.Context, user domain.UserGl) error {
	return a.userGlRepo.Save(ctx, user)
}
