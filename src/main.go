package main

import (
	"web"
	// "fmt"
	// "log"
	// "db"
)

func main() {
	web.Start()
	/*
	db, err := db.Open("E:\\BackUp\\habrahabr_bot\\users.db")
	if err != nil {
		log.Fatal(err)
	}

	elements, err := db.GetCMD()
	for _, elem := range elements {
		fmt.Println(elem.T, "\t", elem.Key, "\t", elem.Value)
	}
	fmt.Print("\n\n\n")
	elements, err = db.Next("users")
	for _, elem := range elements {
		fmt.Println(elem.T, "\t", elem.Key, "\t", elem.Value)
	}
	fmt.Print("\n\n\n")
	elements, err = db.Next("170613028")
	for _, elem := range elements {
		fmt.Println(elem.T, "\t", elem.Key, "\t", elem.Value)
	}
	fmt.Print("\n\n\n")
	db.Back()
	elements, err = db.Back()
	for _, elem := range elements {
		fmt.Println(elem.T, "\t", elem.Key, "\t", elem.Value)
	}
	*/
}