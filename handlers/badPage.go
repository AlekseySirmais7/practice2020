package handlers

import (
	"bytes"
	"database/sql"
	route "github.com/qiangxue/fasthttp-routing"
	"html/template"
	"../serverUtils"
	// необходимо закоментировать ОДНУ строку ниже в импорте
	_ "html/template" // закомментировать эту строку чтобы была xss
	//_"text/template" // закомментировать эту строку чтобы НЕ было xss
)

var PostBadPageTmpl = `
<html><body>
	&lt;script&gt;alert(document.cookie)&lt;/script&gt;
	<br />
	<h2>В serverUtils::CreateCookie(): authCookie.SetHTTPOnly(false) // необходимо установить в false для доступа к cookie</h2>
	<br />
	<h1> Названия постов </h1>

    {{range .Posts}}
		<div style="border: 1px solid black; padding: 5px; margin: 10px;">
			{{.}}
		</div>
    {{end}}
</body></html>`

func (api *ProfileHandler) GetBadPage(ctx *route.Context, myDb *sql.DB) {

	sqlRequest := `SELECT title from Posts 
				   ORDER BY created;`
	rowsSelectedPost, err := myDb.Query(sqlRequest)
	if err != nil {
		panic(err) // bad in real service
	}
	defer rowsSelectedPost.Close()

	postsTitles := []string{}
	for rowsSelectedPost.Next() {
		onePostTitle := ""
		err := rowsSelectedPost.Scan(&onePostTitle)
		if err != nil {
			panic("не получислось получить пост")
		}
		postsTitles = append(postsTitles, onePostTitle)
	}

	tmpl := template.New("main")
	tmpl, _ = tmpl.Parse(PostBadPageTmpl)

	buf := new(bytes.Buffer)

	PostTitles := []template.HTML{}
	for _, v := range postsTitles {
		PostTitles = append(PostTitles, template.HTML(v))
	}
	tmpl.Execute(buf, struct {
		Posts []template.HTML
	}{
		Posts: PostTitles,
	})

	authCookie := serverUtils.CreateCookie("sessionBadPage", "Secret cookie hash", 100000)
	ctx.Response.Header.SetCookie(authCookie)

	ctx.SetContentType("text/html; charset=utf-8")
	ctx.SetStatusCode(200)
	ctx.SetBody([]byte(buf.String()))
}
