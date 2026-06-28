package main

import (
	"os"
	"strconv"

	"github.com/beego/beego/v2/server/web"
	"github.com/kakakikikeke/memo/internal/controller"
	"github.com/kakakikikeke/memo/internal/repository"
	"github.com/kakakikikeke/memo/internal/service"
)

func NewMemoService(repo repository.Repository) *service.MemoService {
	return service.NewMemoService(repo)
}

func main() {
	port := 8080
	if p := os.Getenv("PORT"); p != "" {
		cp, err := strconv.Atoi(p)
		if err != nil {
			panic(err)
		}
		port = cp
	}
	web.BConfig.Listen.HTTPPort = port
	web.BConfig.WebConfig.Session.SessionOn = true
	// for memo
	web.Router("/", new(controller.MainController), "*:ListText")
	web.Router("/insert", new(controller.MainController), "post:SaveText")
	web.Router("/clear", new(controller.MainController), "post:ClearText")
	// for image
	web.Router("/image", new(controller.MainController), "get:ListImage")
	web.Router("/save", new(controller.MainController), "post:SaveImage")
	web.Router("/clear_img", new(controller.MainController), "post:ClearImage")
	// for file
	web.Router("/file", new(controller.MainController), "get:ListFile")
	web.Router("/upload", new(controller.MainController), "post:SaveFile")
	web.Router("/clear_file", new(controller.MainController), "post:ClearFile")
	// for user management
	web.Router("/login", new(controller.MainController), "get:Login")
	web.Router("/check", new(controller.MainController), "post:Check")
	web.Router("/logout", new(controller.MainController), "post:Logout")
	web.Router("/signup", new(controller.MainController), "get:SignUp")
	web.Router("/create", new(controller.MainController), "post:Create")
	web.Router("/signoff", new(controller.MainController), "get:SignOff")
	web.Router("/delete", new(controller.MainController), "post:Delete")
	// for template functions
	web.AddFuncMap("rep", controller.Replace)
	web.AddFuncMap("is_first", controller.IsFirst)
	web.AddFuncMap("is_end", controller.IsEnd)
	web.AddFuncMap("get_file_name", controller.GetFileName)
	web.AddFuncMap("get_content", controller.GetContent)
	web.AddFuncMap("safe", controller.Safe)
	web.AddFuncMap("attr", controller.Attr)
	// for error handling
	web.ErrorController(&controller.ErrorController{})
	web.Run()
}
