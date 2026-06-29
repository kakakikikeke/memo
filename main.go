package main

import (
	"os"
	"strconv"

	"github.com/beego/beego/v2/server/web"
	"github.com/kakakikikeke/memo/internal/controller"
	filectrl "github.com/kakakikikeke/memo/internal/controller/file"
	imagectrl "github.com/kakakikikeke/memo/internal/controller/image"
	textctrl "github.com/kakakikikeke/memo/internal/controller/text"
	userctrl "github.com/kakakikikeke/memo/internal/controller/user"
)

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
	web.Router("/", new(textctrl.Controller), "*:List")
	web.Router("/insert", new(textctrl.Controller), "post:Save")
	web.Router("/clear", new(textctrl.Controller), "post:Clear")
	// for image
	web.Router("/image", new(imagectrl.Controller), "get:List")
	web.Router("/save", new(imagectrl.Controller), "post:Save")
	web.Router("/clear_img", new(imagectrl.Controller), "post:Clear")
	// for file
	web.Router("/file", new(filectrl.Controller), "get:List")
	web.Router("/upload", new(filectrl.Controller), "post:Save")
	web.Router("/clear_file", new(filectrl.Controller), "post:Clear")
	// for user management
	web.Router("/login", new(userctrl.Controller), "get:Login")
	web.Router("/check", new(userctrl.Controller), "post:Check")
	web.Router("/logout", new(userctrl.Controller), "post:Logout")
	web.Router("/signup", new(userctrl.Controller), "get:SignUp")
	web.Router("/create", new(userctrl.Controller), "post:Create")
	web.Router("/signoff", new(userctrl.Controller), "get:SignOff")
	web.Router("/delete", new(userctrl.Controller), "post:Delete")
	// for template functions
	web.AddFuncMap("rep", controller.Replace)
	web.AddFuncMap("is_first", controller.IsFirst)
	web.AddFuncMap("is_end", controller.IsEnd)
	web.AddFuncMap("get_file_name", controller.GetFileName)
	web.AddFuncMap("get_content", controller.GetContent)
	web.AddFuncMap("safe", controller.Safe)
	// for error handling
	web.ErrorController(&controller.ErrorController{})
	web.Run()
}
