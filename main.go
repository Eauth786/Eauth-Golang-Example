package main

import (
	"fmt"
	"os"
	"project/eauth" // Eauth API here
	"time"
)

func main() {
	/* Initialize Eauth (Important) */
	if !eauth.Init() {
		fmt.Println("Failed to initialize Eauth: " + eauth.ErrorMessage)
		time.Sleep(3 * time.Second)
		os.Exit(1)
	}

	eauth.ClearConsole() // clear the console

	fmt.Println("▒█▀▀▀ ░█▀▀█ ▒█░▒█ ▀▀█▀▀ ▒█░▒█ ")
	fmt.Println("▒█▀▀▀ ▒█▄▄█ ▒█░▒█ ░▒█░░ ▒█▀▀█ ")
	fmt.Println("▒█▄▄▄ ▒█░▒█ ░▀▄▄▀ ░▒█░░ ▒█░▒█")
	fmt.Println("[1] Login | [2] Register")
	var option string
	fmt.Print("user@eauth:~$ ")
	fmt.Scanln(&option)

	if option == "1" {
		/* Login (username & password) */
		eauth.ClearConsole() // clear the console
		var username string
		fmt.Print("Username: ")
		fmt.Scanln(&username)

		var password string
		fmt.Print("Password: ")
		fmt.Scanln(&password)

		if eauth.Login(username, password) {
			eauth.ClearConsole() // clear the console
			fmt.Println(eauth.LoggedMessage + "\n")
			fmt.Println("Rank: " + eauth.UserRank)
			fmt.Println("Create Date: " + eauth.RegisterDate)
			fmt.Println("Expire Date: " + eauth.ExpireDate)
			fmt.Println("Hardware ID: " + eauth.HWID)
		} else {
			fmt.Println(eauth.ErrorMessage)
		}
		time.Sleep(3 * time.Second)
		main() // return

	} else if option == "2" {
		/* Register (username & email & password & key) */
		eauth.ClearConsole() // clear the console
		var username string
		fmt.Print("Username: ")
		fmt.Scanln(&username)

		var email string
		fmt.Print("Email: ")
		fmt.Scanln(&email)

		var password string
		fmt.Print("Password: ")
		fmt.Scanln(&password)

		var key string
		fmt.Print("License Key: ")
		fmt.Scanln(&key)

		if eauth.Register(username, email, password, key) {
			eauth.ClearConsole() // clear the console
			fmt.Println(eauth.RegisteredMessage)
		} else {
			fmt.Println(eauth.ErrorMessage)
		}
		time.Sleep(3 * time.Second)
		main() // return
	} else {
		/* None */
		eauth.ClearConsole() // clear the console
		fmt.Println("Invalid option!")
		time.Sleep(3 * time.Second)
		main() // return
	}
}
