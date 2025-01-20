package web

import (
	"context"

	"github.com/odit-bit/sone/internal/monolith"
	"github.com/odit-bit/sone/media/mediahttp"
	"github.com/odit-bit/sone/streaming/streamingpb"
	"github.com/odit-bit/sone/users/gluserpb"
	"github.com/odit-bit/sone/users/userpb"
)

var (
	_googleClientID = "476925020269-3rljlti7sun9pbp7lpbqg04cj3mgqr0k.apps.googleusercontent.com"
)

func StartModule(mono monolith.Monolith) {

	/*below instance act as a service to manage data that out of bounded context of this module (driver).*/
	mc := mediahttp.NewClient(mono.HTTP().Address())

	//dial to (g)rpc server provided by monolith
	conn, err := Dial(context.Background(), mono.RPC().Address())
	if err != nil {
		mono.Logger().Panic(err)
	}
	streamClient := streamingpb.NewLiveStreamClient(conn)
	userClient := userpb.NewUserServiceClient(conn)
	gluserClient := gluserpb.NewGoogleUserServiceClient(conn)

	/* below instance act as a repo or type that manage bounded context of this module (driven) */
	sm := NewSessionManager()

	/*handler use above instance (driven or driver) to serve request */
	// http multiplexer
	mux := mono.Mux()

	//login page renderer
	loginTmpl := NewLoginTemplate(LoginPageArgs{
		AuthCallbackUrl: "/auth",
		RegisterUrl:     "/register",
		GsiUrl:          "/gsi",
		SuccessRedirect: "/stream",
	})
	mux.Get("/login", loginTmpl.HandleFunc)

	//login callback endpoint that will invoke internal api call
	authCB := NewAuthCallback(sm, userClient, gluserClient)
	mux.Post("/auth", authCB.Handle)

	// GSI
	gsi := NewGSIHandler(sm, gluserClient, GSIParam{
		ClientID:        _googleClientID,
		CallbackUrl:     "/auth/gsi",
		SuccessRedirect: "/stream",
	})
	mux.Get("/gsi", gsi.RenderAndHandleAuth)
	mux.Post("/auth/gsi", gsi.HandleGoogle)

	// register handler
	reg := NewRegisterHandler(userClient, sm)
	reg.Handle("/register", "/register/callback", "/login", mux)

	// stream handler
	stream := NewStreamHandler(sm, streamClient, mc)
	stream.Handle("/stream", mux)

}
