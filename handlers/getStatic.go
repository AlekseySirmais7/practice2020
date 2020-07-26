package handlers

import (
	route "github.com/qiangxue/fasthttp-routing"
	"log"
	"../serverUtils"
)

func GetIndex(ctx *route.Context) {
	errSend := serverUtils.SendStaticFile(ctx, "./static/index.html", "text/html")
	if errSend != nil {
		log.Println(" static files: err send (read):", errSend)
	}
}

func GetIndexJs(ctx *route.Context) {
	errSend := serverUtils.SendStaticFile(ctx, "./static/js/index.js", "application/javascript")
	if errSend != nil {
		log.Println(" static files: err send (read):", errSend)
	}
}

func GetNetworkJs(ctx *route.Context) {
	errSend := serverUtils.SendStaticFile(ctx, "./static/js/network.js", "application/javascript")
	if errSend != nil {
		log.Println(" static files: err send (read):", errSend)
	}
}

func GetUserJs(ctx *route.Context) {
	errSend := serverUtils.SendStaticFile(ctx, "./static/js/currentUser.js", "application/javascript")
	if errSend != nil {
		log.Println(" static files: err send (read):", errSend)
	}
}

func GetCss(ctx *route.Context) {
	errSend := serverUtils.SendStaticFile(ctx, "./static/index.css", "text/css")
	if errSend != nil {
		log.Println(" static files: err send (read):", errSend)
	}
}

func GetBootstrapJs(ctx *route.Context) {
	errSend := serverUtils.SendStaticFile(ctx, "./static/bootstrap/bootstrap.min.js", "application/javascript")
	if errSend != nil {
		log.Println(" static files: err send (read):", errSend)
	}
}

func GetBootstrapCss(ctx *route.Context) {
	errSend := serverUtils.SendStaticFile(ctx, "./static/bootstrap/bootstrap.min.css", "text/css")
	if errSend != nil {
		log.Println(" static files: err send (read):", errSend)
	}
}

func GetFavicon(ctx *route.Context) {
	errSend := serverUtils.SendStaticFile(ctx, "./static/favicon.ico", "image/x-icon")
	if errSend != nil {
		log.Println(" static files: err send (read):", errSend)
	}
}

func GetPlusImg(ctx *route.Context) {
	errSend := serverUtils.SendStaticFile(ctx, "./static/plusImg.svg", "image/svg+xml")
	if errSend != nil {
		log.Println(" static files: err send (read):", errSend)
	}
}

func GetCloseImg(ctx *route.Context) {
	errSend := serverUtils.SendStaticFile(ctx, "./static/closeContent.svg", "image/svg+xml")
	if errSend != nil {
		log.Println(" static files: err send (read):", errSend)
	}
}
