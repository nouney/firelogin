package firelogin

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/rs/cors"
	"github.com/skratchdot/open-golang/open"
)

type Firelogin struct {
	// Port to listen to (in case of local HTTP)
	port string
	// Firebase API key
	apiKey string
	// Firebase auth domain
	authDomain string
	// HTML page to serve and to open in a browser to authenticate against Firebase
	authHTML string
	// HTML page to serve and to redirect to if the authentication was successfull
	successHTML string
}

// New creates a new Firelogin instance
// apiKey and authDomain can be found in firebase.
func New(apiKey, authDomain string, opts ...Opt) (*Firelogin, error) {
	ret := &Firelogin{
		apiKey:     apiKey,
		authDomain: authDomain,
	}

	for _, opt := range opts {
		err := opt(ret)
		if err != nil {
			return nil, err
		}
	}

	if ret.port == "" {
		WithPort("8000")(ret)
	}

	ui := NewFirebaseUI("firecli")
	if ret.authHTML == "" {
		WithAuthHTML(ui.AuthHTML())(ret)
	}
	if ret.successHTML == "" {
		WithSuccessHTML(ui.SuccessHTML())(ret)
	}
	return ret, nil
}

// Login logs an user in. This blocks until the user authenticates itself (success or fail)
func (f *Firelogin) Login() (*User, error) {
	// http server base url
	url := "http://localhost:" + f.port

	// auth html
	idxTpl, err := template.New("idxTpl").Parse(f.authHTML)
	if err != nil {
		return nil, errors.Wrap(err, "auth html template")
	}

	// success html
	scsTpl, err := template.New("scsTpl").Parse(f.successHTML)
	if err != nil {
		return nil, errors.Wrap(err, "success html template")
	}

	done := make(chan *User, 1)

	// params for templates
	params := struct {
		APIKey      string
		AuthDomain  string
		URL         string
		SuccessURL  string
		CallbackURL string
	}{f.apiKey, f.authDomain, url, url + "/success", url + "/callback"}

	// create HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := idxTpl.Execute(w, &params)
		if err != nil {
			log.Panic(err)
		}
	})
	// called if auth is successfull
	mux.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
		err := scsTpl.Execute(w, &params)
		if err != nil {
			log.Panic(err)
		}
	})
	// /success should make an ajax call to this endpoint with the firebase user object as payload
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		user := User{}
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		done <- &user
	})

	// execute server in background and kill it before return
	srv := http.Server{Addr: ":" + f.port}
	srv.Handler = cors.Default().Handler(mux)
	defer srv.Shutdown(context.Background())
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Println("Can't start http server:", err.Error())
			done <- nil
		}
	}()

	// don't care of error, the user will manually open the link if needed
	open.Run(url)
	fmt.Println("Your browser has been opened to visit:", url)
	usr := <-done
	return usr, nil
}

// RenewAccessToken renews an access token from a refresh token
func (f *Firelogin) RenewAccessToken(refreshToken string) (*User, error) {
	payload := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
	}
	r, err := http.PostForm("https://securetoken.googleapis.com/v1/token?key="+f.apiKey, payload)
	if err != nil {
		return nil, err
	}

	resp := accessTokenRenewResponse{}
	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	usr, err := f.retrieveUserData(resp.AccessToken)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

// retrieveUserData uses the google API to retrieve the user profile from an access topen
func (f *Firelogin) retrieveUserData(access string) (*User, error) {
	r, err := http.PostForm(
		"https://www.googleapis.com/identitytoolkit/v3/relyingparty/getAccountInfo?key="+f.apiKey,
		url.Values{
			"idToken": {access},
		},
	)
	if err != nil {
		return nil, err
	}

	usr := User{}
	err = json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	return &usr, nil
}
