package main

import (
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"
	"../../handlers"
	"../../serverUtils"
)

// CORS для 5000 порта. Для проверки можно запустить копию сервиса на 5000 порту и обращаться к API на 8080 порт.
func CORS(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

		ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
		ctx.Response.Header.Set("Content-Type", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", `Accept, Content-Type,
		Content-Length, Accept-Encoding, Authorization, x-csrf-token`)
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "http://localhost:5000")

		next(ctx)
	}
}

func main() {

	api := &handlers.ProfileHandler{
		ProfileTable: handlers.NewProfileTable(),
	}

	myDb := serverUtils.OpenConnectDb()

	r := routing.New()

	r.Get("/static/js/index.js", func(c *routing.Context) error {
		handlers.GetIndexJs(c)
		return nil
	})

	r.Get("/static/js/network", func(c *routing.Context) error {
		handlers.GetNetworkJs(c)
		return nil
	})

	r.Get("/static/js/currentUser", func(c *routing.Context) error {
		handlers.GetUserJs(c)
		return nil
	})

	r.Get("/index.css", func(c *routing.Context) error {
		handlers.GetCss(c)
		return nil
	})

	r.Get("/bootstrap/bootstrap.min.css", func(c *routing.Context) error {
		handlers.GetBootstrapCss(c)
		return nil
	})

	r.Get("/bootstrap/bootstrap.min.js", func(c *routing.Context) error {
		handlers.GetBootstrapJs(c)
		return nil
	})

	r.Get("/favicon.ico", func(c *routing.Context) error {
		handlers.GetFavicon(c)
		return nil
	})

	r.Get("/plusImg.svg", func(c *routing.Context) error {
		handlers.GetPlusImg(c)
		return nil
	})

	r.Get("/closeImg.svg", func(c *routing.Context) error {
		handlers.GetCloseImg(c)
		return nil
	})

	r.Post("/api/signup", func(c *routing.Context) error {
		api.SignUp(c, myDb)
		return nil
	})

	r.Post("/api/signin", func(c *routing.Context) error {
		api.SignIn(c, myDb)
		return nil
	})

	r.Post("/api/logout", func(c *routing.Context) error {
		api.LogOut(c)
		return nil
	})

	r.Post("/api/getme", func(c *routing.Context) error {
		api.GetUser(c, myDb)
		return nil
	})

	r.Post("/api/createpost", func(c *routing.Context) error {
		api.CreatePost(c, myDb)
		return nil
	})

	r.Get("/api/getposts", func(c *routing.Context) error {
		api.GetPosts(c, myDb)
		return nil
	})

	r.Get("/getbadpage", func(c *routing.Context) error {
		api.GetBadPage(c, myDb)
		return nil
	})

	r.Get("/*", func(c *routing.Context) error {
		handlers.GetIndex(c)
		return nil
	})

	log.Println("start serving :8080")
	panic(fasthttp.ListenAndServe(":8080", CORS(r.HandleRequest)))

}
