package handlers

import (
	"database/sql"
	route "github.com/qiangxue/fasthttp-routing"
	"../serverUtils"
)

func (api *ProfileHandler) GetUser(ctx *route.Context, myDb *sql.DB) {

	session := string(ctx.Request.Header.Cookie("session"))
	login, ok := api.ProfileTable.Sessions[session]
	answer := userReturnModel{
		Name:  "null",
		Email: "null",
	}
	if !ok {
		serverUtils.SendCommonResponse(answer, ctx, 200) // no session -> send NULL
		return
	}
	rowsCheckUserExist := myDb.QueryRow(`SELECT U.email FROM Users U
 											   WHERE  LOWER(U.login) = LOWER($1)`, login)
	errScanPostId := rowsCheckUserExist.Scan(&answer.Email)
	if errScanPostId != nil {
		serverUtils.SendCommonResponse(nil, ctx, 409)
		return
	}
	answer.Name = login
	answer.CSRF = serverUtils.GetRandom512()
	api.ProfileTable.muTokens.Lock()
	api.ProfileTable.Tokens[answer.CSRF] = answer.Name
	api.ProfileTable.muTokens.Unlock()
	serverUtils.SendCommonResponse(answer, ctx, 200)
}
