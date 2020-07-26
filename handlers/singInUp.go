package handlers

import (
	"database/sql"
	"encoding/json"
	route "github.com/qiangxue/fasthttp-routing"
	"log"
	"../serverUtils"
	"sync"
)

type ProfileTable struct {
	Sessions   map[string]string // [session_hash] login
	Tokens     map[string]string // [token] login
	muSessions sync.RWMutex
	muTokens   sync.RWMutex
}

func NewProfileTable() *ProfileTable {
	return &ProfileTable{
		muSessions: sync.RWMutex{},
		Sessions:   make(map[string]string),
		Tokens:     make(map[string]string),
	}
}

type ProfileHandler struct {
	ProfileTable *ProfileTable
}

type SignUpModel struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userReturnModel struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	CSRF  string `json:"csrf"`
}

func (api *ProfileHandler) SignUp(ctx *route.Context, myDb *sql.DB) {

	body := ctx.PostBody()
	createUserUnmarshal := &SignUpModel{}
	errUnmarshal := json.Unmarshal(body, createUserUnmarshal)
	if errUnmarshal != nil {
		log.Println("ERR unmarshal:", errUnmarshal)
		serverUtils.SendCommonResponse(nil, ctx, 400)
		return
	}

	rowsCheckUserExist, ErrQuery := myDb.Query(`SELECT U.login FROM Users U
 											   WHERE  LOWER(U.email) = LOWER($1) OR
 											   LOWER(U.login) = LOWER($2)`, createUserUnmarshal.Email, createUserUnmarshal.Name)
	if ErrQuery != nil {
		panic(ErrQuery)
	}
	defer rowsCheckUserExist.Close()

	for rowsCheckUserExist.Next() {
		serverUtils.SendCommonResponse(nil, ctx, 409) // не создаем дубликаты пользователей
		log.Println("cancel create profile: same email or login")
		return
	}

	passwordHash := serverUtils.GetDbPassHash(createUserUnmarshal.Password)
	sqlInsertRequest := `INSERT INTO Users (login, email, password) VALUES ($1, $2, $3);`
	_, errExec := myDb.Exec(sqlInsertRequest, createUserUnmarshal.Name, createUserUnmarshal.Email, passwordHash)
	if errExec != nil {
		panic("err insert post :" + errExec.Error())
	}

	// вместо мапы следует использовать быстрое in memory key-value базу
	sessionValue := serverUtils.GetRandom512() + createUserUnmarshal.Name // createUserUnmarshal.Name - UNIQUE (коллизии сессий невозможны)
	api.ProfileTable.muSessions.Lock()
	api.ProfileTable.Sessions[sessionValue] = createUserUnmarshal.Name
	api.ProfileTable.muSessions.Unlock()

	authCookie := serverUtils.CreateCookie("session", sessionValue, 100000)
	ctx.Response.Header.SetCookie(authCookie)

	answerReg := userReturnModel{
		Name:  createUserUnmarshal.Name,
		Email: createUserUnmarshal.Email,
		CSRF:  serverUtils.GetRandom512(),
	}
	api.ProfileTable.muTokens.Lock()
	api.ProfileTable.Tokens[answerReg.CSRF] = answerReg.Name
	api.ProfileTable.muTokens.Unlock()

	serverUtils.SendCommonResponse(answerReg, ctx, 200)
}

type SignInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (api *ProfileHandler) SignIn(ctx *route.Context, myDb *sql.DB) {

	body := ctx.PostBody()
	signInData := new(SignInInput)
	errUnmarshal := json.Unmarshal(body, signInData)
	if errUnmarshal != nil {
		log.Println("SIGN IN: err Unmarshal user body:", errUnmarshal)
		serverUtils.SendCommonResponse(nil, ctx, 400)
	}

	userLogin, errCheck := serverUtils.CheckDbPassHash(signInData.Email, signInData.Password, myDb)
	if errCheck != nil {
		log.Println(" Err CheckDbPassHash():" + errCheck.Error())
		serverUtils.SendCommonResponse(nil, ctx, 403) // Отказ в доступе
		return
	}

	// вместо мапы следует использовать быстрое in memory key-value базу
	sessionValue := serverUtils.GetRandom512() + userLogin
	api.ProfileTable.muSessions.Lock()
	api.ProfileTable.Sessions[sessionValue] = userLogin
	api.ProfileTable.muSessions.Unlock()

	authCookie := serverUtils.CreateCookie("session", sessionValue, 100000)
	ctx.Response.Header.SetCookie(authCookie)

	answerReg := userReturnModel{
		Name:  userLogin,
		Email: signInData.Email,
		CSRF:  serverUtils.GetRandom512(),
	}
	api.ProfileTable.muTokens.Lock()
	api.ProfileTable.Tokens[answerReg.CSRF] = answerReg.Name
	api.ProfileTable.muTokens.Unlock()

	serverUtils.SendCommonResponse(answerReg, ctx, 200)
}

func (api *ProfileHandler) LogOut(ctx *route.Context) {

	_, errCheck := CheckSessionAndCSRF(ctx, api)
	if errCheck != nil {
		log.Println("LogOut(): access denied:", errCheck.Error())
		serverUtils.SendCommonResponse(nil, ctx, 403) // отказано в доступе
		return
	}

	api.ProfileTable.muSessions.Lock()
	delete(api.ProfileTable.Sessions, string(ctx.Request.Header.Cookie("session")))
	api.ProfileTable.muSessions.Unlock()

	authCookie := serverUtils.CreateCookie("session", "NULL", -1)
	ctx.Response.Header.SetCookie(authCookie)
	//ctx.Response.Header.DelCookie("session")
	serverUtils.SendCommonResponse(nil, ctx, 200)
}
