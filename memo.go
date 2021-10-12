package main

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
	"strings"
)

const KEY = "memo"

type Memo struct {
	Msg string `form:"msg"`
}

type User struct {
	Name string `form:"name"`
	Pass string `form:"pass"`
}

type MainController struct {
	web.Controller
}

func NewClient() (client *redis.Client) {
	redisURL := "redis://localhost:6379/0"
	if url := os.Getenv("REDIS_URL"); url != "" {
		redisURL = url
	}
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}
	client = redis.NewClient(opt)
	return client
}

func (mc *MainController) List() {
	key := KEY
	name := mc.GetSession("user")
	if name != nil {
		key = name.(string) + ":" + KEY
		mc.Data["isLogin"] = true
		mc.Data["name"] = name
	}
	logs.Debug(name)
	memos, err := NewClient().LRange(key, 0, -1).Result()
	if err != nil {
		panic(err)
	}
	logs.Debug(memos)
	mc.Layout = "layout.tpl"
	mc.TplName = "list.tpl"
	mc.LayoutSections = make(map[string]string)
	mc.LayoutSections["Header"] = "header.tpl"
	mc.LayoutSections["Scripts"] = "scripts.tpl"
	mc.Data["memos"] = memos
}

func (mc *MainController) Login() {
	mc.Layout = "layout.tpl"
	mc.TplName = "login.tpl"
	mc.LayoutSections = make(map[string]string)
	mc.LayoutSections["Header"] = "header.tpl"
	mc.LayoutSections["Scripts"] = "scripts.tpl"
}

func (mc *MainController) Logout() {
	mc.Data["isLogin"] = false
	mc.Data["name"] = nil
	mc.DelSession("user")
	mc.TplName = "success.tpl"
}

func (mc *MainController) Check() {
	u := User{}
	if err := mc.ParseForm(&u); err != nil {
		panic(err)
	}
	pass, err := NewClient().Get(u.Name).Result()
	if err != nil {
		mc.Ctx.ResponseWriter.WriteHeader(403)
		mc.Data["json"] = map[string]string{"msg": "User does not found."}
		mc.ServeJSON()
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(u.Pass))
	if err != nil {
		mc.Ctx.ResponseWriter.WriteHeader(403)
		mc.Data["json"] = map[string]string{"msg": "Login failed."}
		mc.ServeJSON()
		return
	}
	mc.SetSession("user", u.Name)
	mc.TplName = "success.tpl"
}

func (mc *MainController) Clear() {
	key := KEY
	name := mc.GetSession("user")
	if name != nil {
		key = name.(string) + ":" + KEY
	}
	ret, err := NewClient().Del(key).Result()
	if err != nil {
		panic(err)
	}
	logs.Debug(ret)
	mc.Redirect("/", 302)
}

func (mc *MainController) Insert() {
	m := Memo{}
	if err := mc.ParseForm(&m); err != nil {
		panic(err)
	}
	key := KEY
	name := mc.GetSession("user")
	if name != nil {
		key = name.(string) + ":" + KEY
	}
	size, err := NewClient().LPush(key, m.Msg).Result()
	if err != nil {
		panic(err)
	}
	logs.Debug(size)
	mc.Redirect("/", 302)
}

func replace(str string, from string, to string) (out string) {
	return strings.Replace(str, from, to, -1)
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
	web.Router("/clear", new(MainController), "post:Clear")
	web.Router("/insert", new(MainController), "post:Insert")
	web.Router("/", new(MainController), "*:List")
	web.Router("/login", new(MainController), "get:Login")
	web.Router("/check", new(MainController), "post:Check")
	web.Router("/logout", new(MainController), "post:Logout")
	web.AddFuncMap("rep", replace)
	web.Run()
}
