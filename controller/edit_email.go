package controller

import (
	"backend_test/db"
	U "backend_test/utils"
	"net/http"
	"strings"
)

func EditEmail(w http.ResponseWriter, r *http.Request) {
	var valid = U.HandleBearerToken(w, r)

	if !valid {
		return
	}

	var username = *(U.GetUsernameFromToken(r))
	body := U.GetMap(r)
	email := strings.ToLower(U.CastString(body["email"]))

	if email == "" {
		U.ReturnFailure(w, "Email must not be empty")
		return
	}

	rows, err := db.Instance.Query(`SELECT id FROM "_user" where lower("username") = '` + username + `'`)
	db.CheckError(err)

	defer rows.Close()
	var userId int

	for rows.Next() {
		err = rows.Scan(&userId)
		db.CheckError(err)
	}

	if userId == 0 {
		U.ReturnFailure(w, "User not found")
		return
	}

	updateStmt := `update "_user" set "email"=$1 where "id"=$2`
	_, e := db.Instance.Exec(updateStmt, email, userId)
	db.CheckError(e)

	U.ReturnSuccess(w, map[string]interface{}{
		"id":       userId,
		"username": username,
		"email":    email,
	})

	return
}
