# firelogin

Allows your Firebase users to authenticate in a CLI app.

## Demo

Install the demo:
```shell
$ go install github.com/nouney/firelogin/demo
```

Run it:
```shell
$ $GOPATH/bin/demo
Your browser has been opened to visit: http://localhost:8080

Authentication successfull. Welcome, <your full name>.
```

## Getting started

Install `firelogin`:
```shell
$ go get github.com/nouney/firelogin
```

Copy/paste the code below then run it:

```golang
package main

import (
	"fmt"
	"log"

	"github.com/nouney/firelogin"
)

func main() {
	flogin := firelogin.New(&firelogin.Config{
		APIKey:      "<YOUR FIREBASE API KEY>",
		AuthDomain:  "<YOUR FIREBASE AUTH DOMAIN>",
	})
	// This will block until the user sign in
	user, err := flogin.Login()
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Authentication successfull! Welcome,", user.DisplayName)
}
```

It will open a [FirebaseUI](https://github.com/firebase/firebaseui-web) webpage allowing you to authenticate.

### Customization

#### FirebaseUI

```golang
package main

import (
	"fmt"
	"log"

	"github.com/nouney/firelogin"
)

func main() {
    // no providers = all
	ui := firelogin.NewFirebaseUI(
		"AppName", 
		firelogin.GITHUB_AUTH_PROVIDER_ID, 
		firelogin.GOOGLE_AUTH_PROVIDER_ID
	)
	flogin := firelogin.New(&firelogin.Config{
		APIKey:      "<YOUR FIREBASE API KEY>",
		AuthDomain:  "<YOUR FIREBASE AUTH DOMAIN>",
		URL: "https://your-domain.com/yourpage",
	})
	user, err := flogin.Login()
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Authentication successfull! Welcome,", user.DisplayName)
}
```