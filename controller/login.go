package controller

import (
	"backend_test/db"
	U "backend_test/utils"
	"net/http"
	"strings"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	body := U.GetMap(r)
	username := strings.ToLower(U.CastString(body["username"]))
	inputPassword := U.CastString(body["password"])

	if username == "" {
		U.ReturnFailure(w, "Username must not be empty")
		return
	} else {
		if len(username) < 4 {
			U.ReturnFailure(w, "Username must consists min. 4 character")
			return
		}
	}

	if inputPassword == "" {
		U.ReturnFailure(w, "Password must not be empty")
		return
	} else {
		if len(inputPassword) < 6 {
			U.ReturnFailure(w, "Password must consists min. 6 character")
			return
		}
	}

	rows, err := db.Instance.Query(`SELECT id, password FROM "_user" where lower("username") = '` + username + `'`)
	db.CheckError(err)

	defer rows.Close()
	var userId int
	var password string

	for rows.Next() {
		err = rows.Scan(&userId, &password)
		db.CheckError(err)
	}

	if userId == 0 {
		U.ReturnFailure(w, "User not found")
		return
	}

	if password != inputPassword {
		U.ReturnFailure(w, "Invalid password")
		return
	}

	var signingKey = U.RandStringRunes(12)
	var token = U.GetNewToken(username, signingKey)
	var now = time.Now().UTC()

	deleteStmt := `delete from "_session" where user_id=$1`
	_, e := db.Instance.Exec(deleteStmt, userId)
	db.CheckError(e)

	insertDynStmt := `insert into "_session"("user_id", "created_timestamp", "signing_key") values($1, $2, $3)`
	_, err = db.Instance.Exec(insertDynStmt, userId, now, signingKey)

	db.CheckError(err)

	U.ReturnSuccess(w, map[string]interface{}{
		"token":      token["token"],
		"validUntil": token["validUntil"],
	})

	return
}
