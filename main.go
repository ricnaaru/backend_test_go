package main

import (
	R "backend_test/controller"
	"backend_test/db"
	"net/http"
)

func main() {
	db.Init()
	http.HandleFunc("/register", R.Register)
	http.HandleFunc("/login", R.Login)
	http.HandleFunc("/get-info", R.GetUserInformation)
	http.HandleFunc("/edit-email", R.EditEmail)
	http.ListenAndServe(":8880", nil)
}
