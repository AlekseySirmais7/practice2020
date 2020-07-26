package handlers

import (
	"bytes"
	"errors"
	route "github.com/qiangxue/fasthttp-routing"
	"../serverUtils"
)

func CheckSessionAndCSRF(ctx *route.Context, api *ProfileHandler) (string, error) {

	session := string(ctx.Request.Header.Cookie("session"))
	csrfToken := string(ctx.Request.Header.Peek("X-CSRF-Token"))

	loginFromCSRF, existToken := api.ProfileTable.Tokens[csrfToken]
	if !existToken {
		return "", errors.New("no token")
	}

	loginFromSession, existSession := api.ProfileTable.Sessions[session]
	if !existSession {
		authCookie := serverUtils.CreateCookie("session", "NULL", -1)
		ctx.Response.Header.SetCookie(authCookie)
		return "", errors.New("no session")
	}

	if !bytes.Equal([]byte(loginFromCSRF), []byte(loginFromSession)) {
		return "", errors.New("token user != session user")
	}
	return loginFromSession, nil
}
