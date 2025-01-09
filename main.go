package main

import (
	"fmt"
	"log"
	"project/db"
	"project/models"
	"project/repository"
)

func main() {
	// Connect to the database
	dbc := db.Setup()
	defer func() {
		err := dbc.Close()
		if err != nil {
			log.Fatalf("Failed to close the database: %v", err)
		}
	}()

	// Example usage
	user := models.User{
		ID:    32,
		Name:  "John Doe",
		Email: "john.doe@gmail.com",
	}

	err := repository.CreateUser(dbc, user)
	if err != nil {
		log.Fatal(err)
	}

	var users []models.User
	users = repository.GetAllUsers(dbc)

	fmt.Println(users)
	// TODO: Add application logic here
}
