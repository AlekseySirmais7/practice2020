package handlers

import (
	"database/sql"
	"encoding/json"
	route "github.com/qiangxue/fasthttp-routing"
	"log"
	"../serverUtils"
	"time"
)

type Post struct {
	id     uint
	author string
	title  string
	text   string
	date   string
}

type CreatePostModel struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func (api *ProfileHandler) CreatePost(ctx *route.Context, myDb *sql.DB) {

	body := ctx.PostBody()
	createPostUnmarshal := &CreatePostModel{}
	errUnmarshal := json.Unmarshal(body, createPostUnmarshal)
	if errUnmarshal != nil {
		log.Println("ERR unmarshal:", errUnmarshal)
		serverUtils.SendCommonResponse(nil, ctx, 400)
		return
	}

	login, errCheck := CheckSessionAndCSRF(ctx, api)
	if errCheck != nil {
		log.Println("CreatePost(): access denied:", errCheck.Error())
		serverUtils.SendCommonResponse(nil, ctx, 403) // отказано в доступе
		return
	}

	rowsPostExist, errCheckTitle := myDb.Query(`SELECT P.id FROM Posts P
 										  WHERE  LOWER(P.title) = LOWER($1);`, createPostUnmarshal.Title)
	if errCheckTitle != nil {
		panic(errCheckTitle)
	}
	defer rowsPostExist.Close()
	for rowsPostExist.Next() {
		serverUtils.SendCommonResponse(nil, ctx, 409) // не создаем дубликаты тем новостей
		log.Println("cancel create post: same title")
		return
	}

	insertTime := time.Now()
	sqlInsertRequest := `INSERT INTO Posts (author, title, ttext, created ) VALUES ($1, $2, $3, $4);`
	_, errExec := myDb.Exec(sqlInsertRequest, login, createPostUnmarshal.Title, createPostUnmarshal.Text, insertTime)
	if errExec != nil {
		panic("err insert post :" + errExec.Error()) // в реальных проектах без паники
	}

	newPost := PostModel{Title: createPostUnmarshal.Title,
		Text:   createPostUnmarshal.Text,
		Author: login,
		Date:   insertTime.Format("2006-01-02T15:04:05")}
	serverUtils.SendCommonResponse(newPost, ctx, 200)
}

type PostModel struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	Author string `json:"author"`
	Date   string `json:"date"`
}

func (api *ProfileHandler) GetPosts(ctx *route.Context, myDb *sql.DB) {

	_, errCheck := CheckSessionAndCSRF(ctx, api)
	if errCheck != nil {
		log.Println("GetPosts(): access denied:", errCheck.Error())
		serverUtils.SendCommonResponse(nil, ctx, 403) // отказано в доступе
		return
	}

	sqlRequest := `SELECT author, created, title, ttext from Posts 
				   ORDER BY created;`
	rowsSelectedPost, err := myDb.Query(sqlRequest)
	if err != nil {
		panic(err)
	}
	defer rowsSelectedPost.Close()

	selectedPostArr := []PostModel{}
	for rowsSelectedPost.Next() {
		onePost := PostModel{}
		err := rowsSelectedPost.Scan(&onePost.Author, &onePost.Date, &onePost.Title, &onePost.Text)
		if err != nil {
			panic("не получислось получить пост rowsSelectedPost.Scan")
		}
		selectedPostArr = append(selectedPostArr, onePost)
	}
	serverUtils.SendCommonResponse(selectedPostArr, ctx, 200)
}
