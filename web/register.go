package web

import (
	"io"
	"log"
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/odit-bit/sone/users/userpb"
	"github.com/odit-bit/sone/web/internal/session"
)

// REGISTER

// RENDERER/TEMPLATE

type RegisterPageArgs struct {
	Api string // url to default authentication endpoint api
}

var RegisterPageHTML = `
<!DOCTYPE html>
<html>
<body>

	<h1>Register</h1>

	<form method=post action={{.Api}}>
		<label for="fname">username:</label>
		<input type="text" id="fname" name="fname"><br><br>

		<label for="fpass">password:</label>
		<input type="password" id="fpass" name="fpass"><br><br>


		<input type="submit" value="Submit">
	</form>

</body>
</html>
`

type RegisterTemplate struct {
	data RegisterPageArgs
	tmpl *template.Template
}

func NewRegisterTemplate(ApiEndpoint string) *RegisterTemplate {
	t := template.Must(template.New("register_page").Parse(RegisterPageHTML))
	return &RegisterTemplate{
		data: RegisterPageArgs{
			Api: ApiEndpoint,
		},
		tmpl: t,
	}
}

func (r *RegisterTemplate) Render(w io.Writer, data any) error {
	if data == nil {
		data = r.data
	}
	return r.tmpl.Execute(w, data)
}

//////////////////////////////////////////////////////////////////

type RegisterHandler struct {
	renderer Renderer
	repo     userpb.UserServiceClient
	sm       session.Manager

	// redirectUrl string
}

func NewRegisterHandler(userAPI userpb.UserServiceClient, sm session.Manager) *RegisterHandler {
	t := NewRegisterTemplate("")
	return &RegisterHandler{
		renderer: t,
		repo:     userAPI,
		sm:       sm,
	}
}

func (h *RegisterHandler) Handle(path, registerCallbackUrl, redirectUrl string, mux *chi.Mux) {
	mux.Get(path, h.render(registerCallbackUrl))
	mux.Post(registerCallbackUrl, h.register(redirectUrl))
}

func (h *RegisterHandler) render(registerCallbackUrl string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.renderer.Render(w, RegisterPageArgs{Api: registerCallbackUrl})
		if err != nil {
			log.Println("register Handler failed render html:", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}

func (h *RegisterHandler) register(redirectUrl string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uname := r.FormValue("fname")
		pass := r.FormValue("fpass")
		resp, err := h.repo.RegisterUser(r.Context(), &userpb.RegisterUserRequest{
			Username: uname,
			Password: pass,
		})

		if err != nil {
			log.Println("register handler failed register new user:", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		sess, err := h.sm.Load(r)
		if err != nil {
			log.Println("register handler failed load session:", err)
			http.Error(w, err.Error(), 500)
			return
		}

		sess.Put("userid", resp.GetId())
		err = h.sm.Save(sess, w)
		if err != nil {
			log.Println("register handler failed load session:", err)
			http.Error(w, err.Error(), 500)
			return
		}

		log.Println("register handler succeed register new user, id:", resp.Id)
		http.Redirect(w, r, redirectUrl, http.StatusMovedPermanently)
	}

}
