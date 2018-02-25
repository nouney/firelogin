package firelogin

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/rs/cors"
	"github.com/skratchdot/open-golang/open"
)

type Firelogin struct {
	Config

	user *User
}

type Config struct {
	// Firebase API Key
	APIKey string
	// Firebase Auth domain
	AuthDomain string
	// Port to listen to
	Port string
	// URL to open to authenticate
	URL string
	// HTML for auth page
	AuthHTML string
	// HTML for success page
	SuccessHTML string
}

type User struct {
	UID           string `json:"uid"`
	DisplayName   string `json:"displayName"`
	PhotoURL      string `json:"photoURL"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"emailVerified"`
	IsAnonymous   bool   `json:"isAnonymous"`
	ProviderData  []struct {
		UID         string `json:"uid"`
		DisplayName string `json:"displayName"`
		PhotoURL    string `json:"photoURL"`
		Email       string `json:"email"`
		ProviderID  string `json:"providerId"`
	} `json:"providerData"`
	APIKey          string `json:"apiKey"`
	AppName         string `json:"appName"`
	AuthDomain      string `json:"authDomain"`
	StsTokenManager struct {
		APIKey         string `json:"apiKey"`
		RefreshToken   string `json:"refreshToken"`
		AccessToken    string `json:"accessToken"`
		ExpirationTime int64  `json:"expirationTime"`
	} `json:"stsTokenManager"`
	LastLoginAt string `json:"lastLoginAt"`
	CreatedAt   string `json:"createdAt"`
}

type accessTokenRenewResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	UserID       string `json:"user_id"`
	ProjectID    string `json:"project_id"`
}

func New(c *Config) *Firelogin {
	ret := Firelogin{*c, nil}
	if ret.Port == "" {
		ret.Port = "8080"
	}
	if ret.URL == "" {
		ui := NewFirebaseUI("firecli")
		if ret.AuthHTML == "" {
			ret.AuthHTML = ui.AuthHTML()
		}
		if ret.SuccessHTML == "" {
			ret.SuccessHTML = ui.SuccessHTML()
		}
	}
	return &ret
}

func (f *Firelogin) Login() (*User, error) {
	done := make(chan struct{}, 1)
	url, setHandlers, err := f.getConf()
	if err != nil {
		return nil, err
	}
	srv := f.startHTTP(done, setHandlers)
	defer srv.Shutdown(context.Background())
	fmt.Println("Your browser has been opened to visit:", url)
	err = open.Run(url)
	if err != nil {
		log.Panic(err)
	}
	<-done
	return f.GetUser(), nil
}

func (f *Firelogin) RenewAccessToken(refreshToken string) (*User, error) {
	payload := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
	}
	r, err := http.PostForm("https://securetoken.googleapis.com/v1/token?key="+f.APIKey, payload)
	if err != nil {
		return nil, err
	}
	resp := accessTokenRenewResponse{}
	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	f.user.StsTokenManager.AccessToken = resp.AccessToken
	return f.retrieveUserData(resp.AccessToken)
}

func (f Firelogin) GetUser() *User {
	return f.user
}

type handlersSetter func(*http.ServeMux)

func (f Firelogin) getConf() (url string, setHandlers handlersSetter, err error) {
	// online
	if f.URL != "" {
		return f.URL, nil, nil
	}
	// embedded
	url = "http://localhost:" + f.Port
	idxTpl, err := template.New("idxTpl").Parse(f.AuthHTML)
	if err != nil {
		return "", nil, err
	}
	scsTpl, err := template.New("scsTpl").Parse(f.SuccessHTML)
	if err != nil {
		return "", nil, err
	}
	params := struct {
		APIKey      string
		AuthDomain  string
		URL         string
		SuccessURL  string
		CallbackURL string
	}{f.APIKey, f.AuthDomain, url, url + "/success", url + "/callback"}
	setHandlers = func(mux *http.ServeMux) {
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			err := idxTpl.Execute(w, &params)
			if err != nil {
				log.Panic(err)
			}
		})
		mux.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
			err := scsTpl.Execute(w, &params)
			if err != nil {
				log.Panic(err)
			}
		})
	}
	return url, setHandlers, nil
}

func (f *Firelogin) startHTTP(done chan<- struct{}, setHandlers func(*http.ServeMux)) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		user := User{}
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		f.user = &user
		done <- struct{}{}
	})
	if setHandlers != nil {
		setHandlers(mux)
	}
	srv := http.Server{Addr: ":" + f.Port}
	srv.Handler = cors.Default().Handler(mux)
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Println("Can't start http server:", err.Error())
			done <- struct{}{}
		}
	}()
	return &srv
}

func (f *Firelogin) retrieveUserData(accessToken string) (*User, error) {
	r, err := http.PostForm(
		"https://www.googleapis.com/identitytoolkit/v3/relyingparty/getAccountInfo?key="+f.APIKey,
		url.Values{
			"idToken": {accessToken},
		},
	)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(r.Body).Decode(f.user)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return f.user, nil
}
