package controller

import (
	"backend_test/db"
	U "backend_test/utils"
	"net/http"
)

func GetUserInformation(w http.ResponseWriter, r *http.Request) {
	var valid = U.HandleBearerToken(w, r)

	if !valid {
		return
	}

	var username = *(U.GetUsernameFromToken(r))

	rows, err := db.Instance.Query(`SELECT id, email FROM "_user" where lower("username") = '` + username + `'`)
	db.CheckError(err)

	defer rows.Close()
	var userId int
	var email string

	for rows.Next() {
		err = rows.Scan(&userId, &email)
		db.CheckError(err)
	}

	U.ReturnSuccess(w, map[string]interface{}{
		"id":    userId,
		"email": email,
	})

	return
}
