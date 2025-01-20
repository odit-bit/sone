package web

import (
	"log"
	"net/http"
	"text/template"

	"github.com/odit-bit/sone/users/gluserpb"
	"github.com/odit-bit/sone/web/internal/session"
	"google.golang.org/api/idtoken"
)

var GSIHTML = `
<!DOCTYPE html>
<html>
<body>
	<script src="https://accounts.google.com/gsi/client" async></script> 	
	<script>
		document.addEventListener('htmx:beforeSwap',  function (evt) {
		});

		// Handle successful responses
    	document.body.addEventListener('htmx:afterRequest', function (evt) {
        if (evt.detail.xhr.status === 200) { // Check for success status
            window.location.href = "{{.SuccessRedirect}}"; // Redirect to "/"
        }
    });
	</script>

	<div id="login">
		<div>
			<p>Let's sign in with Google:</p>
			<div
				id="g_id_onload"
				data-client_id="{{.ClientID}}"
				data-login_uri="{{.CallbackUrl}}"
				data-context="signin"
				data-ux_mode="popup"
				data-auto_prompt="false"
				>

			</div>
			<div
				class="g_id_signin"
				data-type="standard"
				data-shape="rectangular"
				data-theme="filled_blue"
				data-text="sign_in_with"
				data-size="large"
				data-logo_alignment="left">
			</div>
		</div>
	</div>
	

	

</body>
</html>
`

type GSIParam struct {
	ClientID        string
	CallbackUrl     string
	SuccessRedirect string
}

type GSIHandler struct {
	sm        session.Manager
	data      GSIParam
	gluserAPI gluserpb.GoogleUserServiceClient
	tmpl      *template.Template
}

func NewGSIHandler(sm session.Manager, gluserAPI gluserpb.GoogleUserServiceClient, param GSIParam) *GSIHandler {

	return &GSIHandler{
		tmpl:      template.Must(template.New("GSIPage").Parse(GSIHTML)),
		data:      param,
		gluserAPI: gluserAPI,
		sm:        sm,
	}
}

func (t *GSIHandler) RenderAndHandleAuth(w http.ResponseWriter, r *http.Request) {
	if err := t.tmpl.Execute(w, t.data); err != nil {
		log.Println("GSI handler failed handle auth", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (h *GSIHandler) HandleGoogle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if err := r.ParseForm(); err != nil {
		log.Println("Auth Callback Err:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//ver crsf token
	token, err := r.Cookie("g_csrf_token")
	if err != nil {
		log.Println("Auth Callback read cookie:", err)
		http.Error(w, "no token found", http.StatusBadRequest)
		return
	}

	bodyToken := r.FormValue("g_csrf_token")
	if token.Value != bodyToken {
		http.Error(w, "token missmatch", http.StatusBadRequest)
		return
	}

	//verify id token
	ctx := r.Context()
	validator, err := idtoken.NewValidator(r.Context())
	if err != nil {
		panic(err)
	}

	//google id token
	credential := r.FormValue("credential")
	payload, err := validator.Validate(ctx, credential, _googleClientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claim := struct {
		id         string
		name       string
		email      string
		isVerified bool
	}{}

	// check claim
	claim.id = payload.Claims["sub"].(string)
	if claim.id == "" {
		log.Println("web auth handler error: google id return from google is not string type ,THIS IS BUG !!")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	claim.name = payload.Claims["name"].(string)
	claim.email = payload.Claims["email"].(string)
	claim.isVerified = payload.Claims["email_verified"].(bool)

	//check existing user
	resp, err := h.gluserAPI.Get(ctx, &gluserpb.GetRequest{
		Id: claim.id,
	})
	if err != nil {
		log.Println("web auth handler failed fetch google user from repo:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	exist := resp.IsFound
	modified := false
	if exist {
		log.Printf("web auth handler: user id found %s", resp.Id)
		if resp.Email != claim.email {
			modified = true
		}

		if resp.Name != claim.name {
			modified = true
		}

		if resp.IsEmailVerified != claim.isVerified {
			modified = true
		}

	}

	if !exist || modified {
		log.Println("web auth handler: update new or existed user ", claim.id)
		_, err := h.gluserAPI.Save(ctx, &gluserpb.SaveRequest{
			Id:              claim.id,
			Name:            claim.name,
			Email:           claim.email,
			IsEmailVerified: claim.isVerified,
		})

		if err != nil {
			log.Println("web auth handler failed to save user:", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	} else {
		// user exist and not modified
	}
	sess, err := h.sm.Load(r)
	if err != nil {
		log.Println("GSI handler failed load session", err)
		http.Error(w, "invalid session", http.StatusBadRequest)
		return
	}

	sess.SetLogin(true)
	sess.SetName(resp.Name)
	err = h.sm.Save(sess, w)
	if err != nil {
		log.Println("GSI handler failed save session", err)
		http.Error(w, "invalid session", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, h.data.SuccessRedirect, http.StatusMovedPermanently)
}
