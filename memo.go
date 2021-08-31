package main

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/go-redis/redis"
	"os"
	"strconv"
	"strings"
)

const KEY = "memo"

type Memo struct {
	Msg string `form:"msg"`
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
	memos, err := NewClient().LRange(KEY, 0, -1).Result()
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

func (mc *MainController) Clear() {
	ret, err := NewClient().FlushAll().Result()
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
	size, err := NewClient().LPush(KEY, m.Msg).Result()
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
	web.Router("/clear", new(MainController), "post:Clear")
	web.Router("/insert", new(MainController), "post:Insert")
	web.Router("/*", new(MainController), "*:List")
	web.AddFuncMap("rep", replace)
	web.Run()
}
