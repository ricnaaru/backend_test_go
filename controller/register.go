package controller

import (
	"backend_test/db"
	U "backend_test/utils"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func Register(w http.ResponseWriter, r *http.Request) {
	body := U.GetMap(r)
	username := strings.ToLower(U.CastString(body["username"]))
	password := U.CastString(body["password"])
	email := strings.ToLower(U.CastString(body["email"]))

	if username == "" {
		U.ReturnFailure(w, "Username must not be empty")
		return
	} else {
		if len(username) < 4 {
			U.ReturnFailure(w, "Username must consists min. 4 character")
			return
		}
	}

	if password == "" {
		U.ReturnFailure(w, "Password must not be empty")
		return
	} else {
		if len(password) < 6 {
			U.ReturnFailure(w, "Password must consists min. 6 character")
			return
		}
	}

	if email == "" {
		U.ReturnFailure(w, "Email must not be empty")
		return
	}

	rows, err := db.Instance.Query(`SELECT Count(1) FROM "_user" where lower("username") = '` + username + `'`)
	db.CheckError(err)

	defer rows.Close()
	var count int

	for rows.Next() {

		err = rows.Scan(&count)
		db.CheckError(err)

		fmt.Println("count => " + strconv.Itoa(count))
	}

	if count > 0 {
		U.ReturnFailure(w, "Username has been taken")
		return
	}

	insertDynStmt := `insert into "_user"("username", "password", "email") values($1, $2, $3)`
	_, err = db.Instance.Exec(insertDynStmt, username, password, email)

	db.CheckError(err)

	U.ReturnSuccess(w, map[string]interface{}{
		"message": "Done!",
	})

	return
}
