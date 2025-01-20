package web

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"text/template"

	"github.com/odit-bit/sone/users/gluserpb"
	"github.com/odit-bit/sone/users/userpb"
	"github.com/odit-bit/sone/web/internal/session"
)

// <script>
//         document.getElementById('loginForm').addEventListener('submit', async function(event) {
//             event.preventDefault(); // Prevent the default form submission

//             const form = event.target;
//             const formData = new FormData(form);

//             try {
//                 const response = await fetch(form.action, {
//                     method: form.method,
//                     body: formData,
//                 });

//                 if (response.ok) {
//                     // Redirect to the home page on success
//                     window.location.href = {{.SuccessRedirect}};
//                 } else {
//                     // Handle authentication failure (e.g., show an error message)
//                     alert("Authentication failed. Please try again.");
// 					}
//             } catch (error) {
//                 console.error("Error during authentication:", error);
//                 alert("An error occurred. Please try again later.");
//             }
//         });
//     </script>
/// LOGIN

type LoginPageArgs struct {
	AuthCallbackUrl string // url to default authentication endpoint api
	// GoogleClientID    string
	// GoogleCallbackUrl string
	RegisterUrl     string
	SuccessRedirect string
	GsiUrl          string
}

var LoginPageHTML = `
<!DOCTYPE html>
<html>
<body>
	<script src="https://accounts.google.com/gsi/client" async></script> 
	<script src="https://unpkg.com/htmx.org@2.0.4"></script>
	
	<script>
		document.addEventListener('htmx:beforeSwap',  function (evt) {
			const xhr = evt.detail.xhr;
			if (evt.detail.xhr.status === 422) {
			  	const errorText = evt.detail.xhr.responseText; // Get the plain text response
            	const errorDiv = document.getElementById('errors');
				// Set the error text inside the #errors div
				errorDiv.innerHTML = errorText;

				// Prevent HTMX from swapping the main target (form)
            	evt.detail.shouldSwap = false;
			}
		
		});

		// Handle successful responses
    	document.body.addEventListener('htmx:afterRequest', function (evt) {
        if (evt.detail.xhr.status === 200) { // Check for success status
            window.location.href = "{{.SuccessRedirect}}"; // Redirect to "/"
        }
    });
	</script>

	<div id="login">
		<h1>Welcome to this web app!</h1>
		<form  hx-post="{{.AuthCallbackUrl}}" enctype="multipart/form-data" hx-swap="outerHTML" hx-target="#login">
			<label for="fname">username:</label>
			<input type="text" id="fname" name="fname"><br><br>

			<label for="fpass">password:</label>
			<input type="password" id="fpass" name="fpass"><br><br>

			<div id="errors"></div>
			<input type="submit" value="Submit">
		</form> 

		<a href={{.RegisterUrl}}>
			<button>sign-up</button>
		</a>

		<a href={{.GsiUrl}}>
			<button>Google</button>
		</a>

		<p>----</p>

	
		
	</div>
	

	

</body>
</html>
`

type LoginPageTemplate struct {
	data LoginPageArgs
	tmpl *template.Template
}

func NewLoginTemplate(param LoginPageArgs) *LoginPageTemplate {
	tmpl := template.Must(template.New("login_page").Parse(LoginPageHTML))

	return &LoginPageTemplate{
		data: param,
		tmpl: tmpl,
	}
}

func (l *LoginPageTemplate) Render(w io.Writer, data any) error {
	if data == nil {
		data = l.data
	}
	return l.tmpl.Execute(w, data)
}

func (l *LoginPageTemplate) HandleFunc(w http.ResponseWriter, r *http.Request) {
	if err := l.tmpl.Execute(w, l.data); err != nil {
		log.Println("loginTmpl failed render:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

/// Auth Failure

var AuthFailureHTML = `
<html>
<body>
	<p>wrong username or password </p>
</body>
</html>
`

type AuthFailureArgs struct {
	LoginPageUrl string
	// Placeholder  string
}

type AuthFailureTemplate struct {
	tmpl *template.Template
}

func NewAuthFailureTemplate() *AuthFailureTemplate {
	tmpl := template.Must(template.New("authFailure").Parse(AuthFailureHTML))
	return &AuthFailureTemplate{
		tmpl: tmpl,
	}
}

func (t *AuthFailureTemplate) Render(w io.Writer, data any) error {
	if data == nil {
		return fmt.Errorf("auth failure template, data cannot be nil this is a bug")
	}
	args, ok := data.(AuthFailureArgs)
	if !ok {
		return fmt.Errorf("auth failure template, wrong data type got %T expected %T ", data, AuthFailureArgs{})
	}

	return t.tmpl.Execute(w, args)
}

type AuthCallbackHandler struct {
	sm        session.Manager
	userAPI   userpb.UserServiceClient
	gluserAPI gluserpb.GoogleUserServiceClient
}

func NewAuthCallback(sm session.Manager, userAPI userpb.UserServiceClient, glUserAPI gluserpb.GoogleUserServiceClient) *AuthCallbackHandler {
	return &AuthCallbackHandler{
		sm:        sm,
		userAPI:   userAPI,
		gluserAPI: glUserAPI,
	}
}

func (h *AuthCallbackHandler) Handle(w http.ResponseWriter, r *http.Request) {

	uname := r.FormValue("fname")
	pass := r.FormValue("fpass")

	sess, err := h.sm.Load(r)
	if err != nil {
		log.Println("auth handler failed load session :", err)
		http.Error(w, "server error", http.StatusUnprocessableEntity)
	}

	resp, err := h.userAPI.AuthenticateUser(r.Context(), &userpb.AuthUserRequest{
		Username: uname,
		Password: pass,
	})
	if err != nil {
		log.Println("auth handler failed authentication:", err)
		http.Error(w, "wrong username or password", http.StatusUnprocessableEntity)
		return
	}
	sess.Put("token", resp.Token)
	sess.SetLogin(true)
	sess.SetName(uname)

	if err := h.sm.Save(sess, w); err != nil {
		log.Println("auth handler failed save session :", err)
		http.Error(w, "server error", http.StatusUnprocessableEntity)
		return
	}

}
