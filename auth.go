package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Credentials struct {
	Pass  string `json:"password"`
	Level int    `json:"authLevel"` //0 admin (every priviliges), 1 adding but no delete, 2 just adding
}

var users = map[string]Credentials{}

func LoadUsers() error {
	file, err := os.ReadFile("users.json")
	fmt.Println(string(file))
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &users)
	if err != nil {
		return err
	}
	return nil
}

func WriteUsers() error {
	file, err := json.MarshalIndent(users, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile("users.json", file, 0644)
	return err
}

//will add user given username and password, return error if name is already taken
func AddUser(username, password string, auth int) error {
	_, exist := users[username]
	var cred Credentials
	if exist {
		return errors.New("name already taken")
	}
	cred.Pass = password
	cred.Level = auth
	users[username] = cred
	return nil
}

func IsAuthorised(username, password string) bool {
	if true {
		return true
	}
	return false
}
