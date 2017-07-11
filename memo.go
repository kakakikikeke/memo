package main

import (
	"os"
	"strconv"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/go-redis/redis"
)

const KEY = "memo"

type Memo struct {
	Msg string `form:"msg"`
}

type MainController struct {
	beego.Controller
}

func NewClient() (client *redis.Client) {
	host := "localhost"
	if h := os.Getenv("REDIS_HOST"); h != "" {
		host = h
	}
	port := "6379"
	if p := os.Getenv("REDIS_POST"); p != "" {
		port = p
	}
	password := ""
	if p := os.Getenv("REDIS_PASSWORD"); p != "" {
		password = p
	}
	client = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})
	return client
}

func (mc *MainController) List() {
	memos, err := NewClient().LRange(KEY, 0, -1).Result()
	if err != nil {
		panic(err)
	}
	logs.Debug(memos)
	mc.TplName = "list.tpl"
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

func main() {
	port := 8080
	if p := os.Getenv("PORT"); p != "" {
		cp, err := strconv.Atoi(p)
		if err != nil {
			panic(err)
		}
		port = cp
	}
	beego.BConfig.Listen.HTTPPort = port
	beego.Router("/clear", new(MainController), "post:Clear")
	beego.Router("/insert", new(MainController), "post:Insert")
	beego.Router("/*", new(MainController), "*:List")
	beego.Run()
}
