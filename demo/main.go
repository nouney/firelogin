package main

import (
	"fmt"
	"log"

	"github.com/nouney/firelogin"
)

func main() {
	ui := firelogin.NewFirebaseUI("Firelogin-demo", firelogin.GITHUB_AUTH_PROVIDER_ID)

	flogin, _ := firelogin.New(
		"AIzaSyCvz8teBEpFBqR6LScQrp_WNcSloZdG8X4",
		"firelogin-demo.firebaseapp.com",
		firelogin.WithAuthHTML(ui.AuthHTML()),
		firelogin.WithSuccessHTML(ui.SuccessHTML()),
	)

	// This will block until the user signs in
	user, err := flogin.Login()
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Authentication successfull! Welcome,", user.DisplayName)
}
