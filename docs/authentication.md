# Authentication

```sh
$ alchemy add authentication.all
```

# Usage

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/struckchure/alchemy-sandbox/dao"
	"github.com/struckchure/alchemy-sandbox/prisma/db"
	"github.com/struckchure/alchemy-sandbox/services"
)

func prettyPrint(i interface{}) {
	s, _ := json.MarshalIndent(i, "", " ")
	fmt.Println(string(s))
}

func prompt(msg string) (*string, error) {
	fmt.Print(msg)

	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return nil, err
	}

	return &input, nil
}

func main() {
	client := db.NewClient()
	client.Connect()
	defer client.Disconnect()

	userDao := dao.NewUserDao(client)
	authenticationService := services.NewAuthenticationService(userDao)

	route, _ := prompt("Route [login or register]: ")

	switch *route {
	case "register":
		firstName, _ := prompt("First name: ")
		lastName, _ := prompt("Last name: ")
		email, _ := prompt("Email: ")
		password, _ := prompt("Password: ")

		registerRes, err := authenticationService.Register(services.RegisterArgs{
			FirstName: firstName,
			LastName:  lastName,
			Email:     *email,
			Password:  *password,
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		prettyPrint(registerRes)
	case "login":
		email, _ := prompt("Email: ")
		password, _ := prompt("Password: ")

		loginRes, err := authenticationService.Login(services.LoginArgs{
			Email:    *email,
			Password: *password,
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		prettyPrint(loginRes)
	default:
		fmt.Println("404 - page not found")
		return
	}
}
```
