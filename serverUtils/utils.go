package serverUtils

import (
	"bytes"
	_ "bytes"
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	_ "github.com/lib/pq"
	route "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/pbkdf2"
	"io/ioutil"
)

func SendCommonResponse(v interface{}, ctx *route.Context, responseStatus int) {
	jsonDataAns, err := json.Marshal(v)
	if err != nil {
		panic(err) // bad in real project
		return
	}
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(responseStatus)
	ctx.SetBody(jsonDataAns)
}

func SendStaticFile(ctx *route.Context, filePath string, contentType string) error {
	fileData, errRead := ioutil.ReadFile(filePath)
	if errRead != nil {
		ctx.SetStatusCode(500)
		ctx.SetBody([]byte("server err"))
		return errRead
	}
	ctx.SetContentType(contentType)
	ctx.SetStatusCode(200)
	ctx.SetBody(fileData)
	return nil
}

func OpenConnectDb() *sql.DB {
	connStr := "user=docker password=docker dbname=myService sslmode=disable"
	myDb, errOpen := sql.Open("postgres", connStr)
	if errOpen != nil {
		panic(errOpen)
	}
	myDb.SetMaxOpenConns(5)
	return myDb
}

func CreateCookie(key string, value string, expire int) *fasthttp.Cookie {
	authCookie := fasthttp.Cookie{}
	authCookie.SetKey(key)
	authCookie.SetValue(value)
	authCookie.SetMaxAge(expire)
	authCookie.SetHTTPOnly(true) // менять тут показа кук через xss
	authCookie.SetSameSite(fasthttp.CookieSameSiteLaxMode)
	return &authCookie
}

func GetDbPassHash(password string) string {
	salt := make([]byte, 32)
	rand.Read(salt)
	passwordHash := pbkdf2.Key([]byte(password), salt, 4096, 512, sha1.New)
	passwordHash = append(salt, passwordHash...) // first 32 bytes is salt
	return string(passwordHash)
}

func CheckDbPassHash(email string, password string, myDb *sql.DB) (string, error) {

	passwordDbHash := ""
	userLogin := "NULL"

	/*
		// may be SQL Injection
		sqlRequest := `SELECT U.password, U.login FROM Users U
					   WHERE  U.email = '` + email + `' AND U.password = '` + password + "'"
		log.Println("may be SQL Injection: " + sqlRequest)
		rowsCheckPass := myDb.QueryRow(sqlRequest)
		errScanPass := rowsCheckPass.Scan(&passwordDbHash, &userLogin)
		if errScanPass != nil {
			return userLogin, errors.New("Not correct")
		}
		return userLogin, nil
	*/

	rowsCheckPass := myDb.QueryRow(`SELECT U.password, U.login FROM Users U
										  WHERE  LOWER(U.email) = LOWER($1)`, email)
	errScanPass := rowsCheckPass.Scan(&passwordDbHash, &userLogin)
	if errScanPass != nil {
		return userLogin, errors.New("No user with input email")
	}

	salt := []byte(passwordDbHash[:32]) // first 32 bytes is salt
	passwordInputHash := pbkdf2.Key([]byte(password), salt, 4096, 512, sha1.New)
	passwordInputHash = append(salt, passwordInputHash...)

	if bytes.Equal(passwordInputHash, []byte(passwordDbHash)) {
		return userLogin, nil
	}
	return userLogin, errors.New("Wrong pass")
}

func GetRandom512() string {
	randomBytes := make([]byte, 512)
	rand.Read(randomBytes)
	return hex.EncodeToString(randomBytes)
}
