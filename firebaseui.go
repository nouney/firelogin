package firelogin

const (
	GOOGLE_AUTH_PROVIDER_ID   = "firebase.auth.GoogleAuthProvider.PROVIDER_ID"
	FACEBOOK_AUTH_PROVIDER_ID = "firebase.auth.FacebookAuthProvider.PROVIDER_ID"
	TWITTER_AUTH_PROVIDER_ID  = "firebase.auth.TwitterAuthProvider.PROVIDER_ID"
	GITHUB_AUTH_PROVIDER_ID   = "firebase.auth.GithubAuthProvider.PROVIDER_ID"
	EMAIL_AUTH_PROVIDER_ID    = "firebase.auth.EmailAuthProvider.PROVIDER_ID"
	PHONE_AUTH_PROVIDER_ID    = "firebase.auth.PhoneAuthProvider.PROVIDER_ID"
)

type FirebaseUI struct {
	name      string
	providers []string
}

func NewFirebaseUI(name string, providers ...string) *FirebaseUI {
	ret := &FirebaseUI{name: name}
	if len(providers) == 0 {
		ret.providers = []string{GOOGLE_AUTH_PROVIDER_ID, FACEBOOK_AUTH_PROVIDER_ID, TWITTER_AUTH_PROVIDER_ID,
			GITHUB_AUTH_PROVIDER_ID, EMAIL_AUTH_PROVIDER_ID, PHONE_AUTH_PROVIDER_ID}
	} else {
		ret.providers = providers
	}
	return ret
}

func (fui FirebaseUI) AuthHTML() string {
	var providers string
	for _, p := range fui.providers {
		providers = providers + p + ",\n"
	}
	return `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>` + fui.name + `</title>
	<script src="https://www.gstatic.com/firebasejs/4.10.0/firebase.js"></script>
	<script>
	// Initialize Firebase
	var config = {
		apiKey: "{{.APIKey}}",
		authDomain: "{{.AuthDomain}}",
	};
	firebase.initializeApp(config);
	</script>
	<script src="https://cdn.firebase.com/libs/firebaseui/2.6.1/firebaseui.js"></script>
	<link type="text/css" rel="stylesheet" href="https://cdn.firebase.com/libs/firebaseui/2.6.1/firebaseui.css" />
	<script type="text/javascript">
	// FirebaseUI config.
	var uiConfig = {
		signInSuccessUrl: '{{.SuccessURL}}',
		signInOptions: [
		` + providers + `
		],
	};
	// Initialize the FirebaseUI Widget using Firebase.
	var ui = new firebaseui.auth.AuthUI(firebase.auth());
	// The start method will wait until the DOM is loaded.
	ui.start('#firebaseui-auth-container', uiConfig);
	</script>
</head>
<body>
	<h1>` + fui.name + `: Sign-in</h1>
	<div id="firebaseui-auth-container"></div>
</body>
</html>
`
}

func (fui FirebaseUI) SuccessHTML() string {
	return `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>` + fui.name + ` - Success</title>
	<script src="https://www.gstatic.com/firebasejs/4.10.0/firebase.js"></script>
	<script>
	// Initialize Firebase
	var config = {
		apiKey: "{{.APIKey}}",
		authDomain: "{{.AuthDomain}}",
	};
	firebase.initializeApp(config);
	</script>
	<script src="https://cdn.firebase.com/libs/firebaseui/2.6.1/firebaseui.js"></script>
	<link type="text/css" rel="stylesheet" href="https://cdn.firebase.com/libs/firebaseui/2.6.1/firebaseui.css" />
	<script type="text/javascript">
	// FirebaseUI config.
	firebase.auth().onAuthStateChanged(function(user) {
		if (user) {
			user.getIdToken().then(function(accessToken) {
			var xhr = new XMLHttpRequest();
			xhr.open('POST', '{{.CallbackURL}}');
			xhr.send(JSON.stringify(user));
			})
		}
	})
	</script>
</head>
<body>
	<!-- The surrounding HTML is left untouched by FirebaseUI.
		Your app may use that space for branding, controls and other customizations.-->
	<h1>` + fui.name + ` - Authentication successfull</h1>
	<div id="firebaseui-auth-container"></div>
</body>
</html>
`
}
