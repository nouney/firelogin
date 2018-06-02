package firelogin

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
