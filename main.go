package main

import (
	"fmt"
	"lesson6/document_store"
	"lesson6/users"
)

func main() {
	storage := document_store.NewStore()

	isCreatedColl, _ := storage.CreateCollection("users", &document_store.CollectionConfig{
		PrimaryKey: "ID",
	})

	if !isCreatedColl {
		fmt.Println("error in creating collection")
	}

	userService, err := users.UserService(storage, "users")

	if err != nil {
		fmt.Println(err)
		return
	}

	_, errU1 := userService.CreateUser("u1", "Kargand")

	if errU1 != nil {
		fmt.Println(errU1)
		return
	}

	_, errU2 := userService.CreateUser("u2", "Beshbarmak")

	if errU2 != nil {
		fmt.Println(errU2)
		return
	}

	listUsers, errList := userService.ListUsers()

	if errList != nil {
		fmt.Println(errList)
		return
	}

	fmt.Println("users:", listUsers)

	user, errUsr2 := userService.GetUser("u2")

	if errUsr2 != nil {
		fmt.Println(errUsr2)
		return
	}

	fmt.Println("Username that was found:", "ID:", user.ID, "Name:", user.Name)
}
