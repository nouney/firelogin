package main

import (
	"fmt"
	"log"

	"github.com/nouney/firelogin"
)

func main() {
	ui := firelogin.NewFirebaseUI("Firelogin-demo", firelogin.GITHUB_AUTH_PROVIDER_ID)
	flogin := firelogin.New(&firelogin.Config{
		APIKey:      "AIzaSyCvz8teBEpFBqR6LScQrp_WNcSloZdG8X4",
		AuthDomain:  "firelogin-demo.firebaseapp.com",
		AuthHTML:    ui.AuthHTML(),
		SuccessHTML: ui.SuccessHTML(),
	})
	// This will block until the user sign in
	user, err := flogin.Login()
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Authentication successfull! Welcome,", user.DisplayName)
}
