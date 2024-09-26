package main

import (
	"fmt"

	"github.com/kodernubie/gocommon/db"
	"github.com/oklog/ulid/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// MODEL -----------------
type User struct {
	ID       string `bson:"_id"`
	Name     string
	Email    string
	Password string
}

// func (o *User) TableName() string {
// 	return "Users"
// }

func (o *User) BeforeCreate() {

	fmt.Println("before insert")
	o.ID = ulid.Make().String()
}

//------------------------------

func main() {

	choice := ""

	for {

		fmt.Println("\r\nMenu :")
		fmt.Println("1. Add User")
		fmt.Println("2. Search User")
		fmt.Println("3. Get User By Id")
		fmt.Println("4. Update User")
		fmt.Println("5. Delete User")
		fmt.Println("6. Quit")
		fmt.Print("Selection : ")

		_, err := fmt.Scanln(&choice)

		if err != nil {
			fmt.Println("Error : ", err)
			return
		}

		switch choice {
		case "1":
			Add()
		case "2":
			Search()
		case "3":
			Get()
		case "4":
			Update()
		case "5":
			Delete()
		case "6":
			return
		default:
			fmt.Println("Invalid choice")
		}
	}
}

func Add() {

	name := ""
	fmt.Print("User Name : ")
	fmt.Scanln(&name)

	email := ""
	fmt.Print("Email : ")
	fmt.Scanln(&email)

	pass := ""
	fmt.Print("Password : ")
	fmt.Scanln(&pass)

	user := &User{
		Name:     name,
		Password: pass,
		Email:    email,
	}

	err := db.Save(user)

	if err != nil {
		fmt.Println("Error create :", err)
		return
	}

	byt, _ := bson.Marshal(user)

	fmt.Println(string(byt))
	fmt.Println("User created with id :", user.ID)
}

func Search() {

	filter := ""
	fmt.Print("Filter : ")
	fmt.Scanln(&filter)

	users := []User{}

	err := db.Find(&users, bson.M{"name": bson.M{"$regex": filter, "$options": "im"}},
		db.FindOption{
			Skip:  1,
			Limit: 2,
			Order: []db.FieldOrder{{"name", "asc"}},
		})

	if err != nil {
		fmt.Println("Error :", err)
		return
	}

	for _, item := range users {

		fmt.Printf("%+v\r\n", item)
	}
}

func Get() {

	id := ""
	fmt.Print("ID : ")
	fmt.Scanln(&id)

	user := User{}

	err := db.FindById(&user, id)

	if err != nil {
		fmt.Println("Error :", err)
		return
	}

	fmt.Printf("%+v\r\n", user)
}

func Update() {

	id := ""
	fmt.Print("ID : ")
	fmt.Scanln(&id)

	user := User{}

	err := db.FindById(&user, id)

	if err != nil {
		fmt.Println("Error :", err)
		return
	}

	fmt.Printf("%+v\r\n", user)

	name := ""
	fmt.Print("User Name : ")
	fmt.Scanln(&name)

	email := ""
	fmt.Print("Email : ")
	fmt.Scanln(&email)

	pass := ""
	fmt.Print("Password : ")
	fmt.Scanln(&pass)

	user.Name = name
	user.Email = email
	user.Password = pass

	err = db.Save(&user)

	if err != nil {
		fmt.Println("Error :", err)
		return
	}

	fmt.Println("User Updated !")
}

func Delete() {

	id := ""
	fmt.Print("ID : ")
	fmt.Scanln(&id)

	user := User{}

	err := db.FindById(&user, id)

	if err != nil {
		fmt.Println("Error :", err)
		return
	}

	db.Delete(user)
	fmt.Println("User deleted !")
}
