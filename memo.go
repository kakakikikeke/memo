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
const IMAGE_KEY = "image"

type Memo struct {
	Msg string `form:"msg"`
}

type User struct {
	Name string `form:"name"`
	Pass string `form:"pass"`
}

type NewUser struct {
	Name  string `form:"name"`
	Pass  string `form:"pass"`
	Pass2 string `form:"pass2"`
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
	cli := NewClient()
	defer cli.Close()
	memos, err := cli.LRange(key, 0, -1).Result()
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

func (mc *MainController) SignUp() {
	mc.Layout = "layout.tpl"
	mc.TplName = "signup.tpl"
	mc.LayoutSections = make(map[string]string)
	mc.LayoutSections["Header"] = "header.tpl"
	mc.LayoutSections["Scripts"] = "scripts.tpl"
}

func (mc *MainController) SignOff() {
	name := mc.GetSession("user")
	if name == nil {
		mc.Redirect("/", 302)
	}
	mc.Data["name"] = name
	mc.Layout = "layout.tpl"
	mc.TplName = "signoff.tpl"
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
	cli := NewClient()
	defer cli.Close()
	pass, err := cli.Get(u.Name).Result()
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

func (mc *MainController) Create() {
	newu := NewUser{}
	if err := mc.ParseForm(&newu); err != nil {
		panic(err)
		return
	}
	cli := NewClient()
	defer cli.Close()
	err := cli.Get(newu.Name).Err()
	if err == nil {
		mc.Ctx.ResponseWriter.WriteHeader(403)
		mc.Data["json"] = map[string]string{"msg": "Specified user already exists."}
		mc.ServeJSON()
		return
	}

	if newu.Pass != newu.Pass2 {
		mc.Ctx.ResponseWriter.WriteHeader(403)
		mc.Data["json"] = map[string]string{"msg": "Does not match password."}
		mc.ServeJSON()
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newu.Pass), bcrypt.DefaultCost)
	err = cli.Set(newu.Name, hash, 0).Err()
	if err != nil {
		panic(err)
		return
	}
	mc.SetSession("user", newu.Name)
	mc.TplName = "success.tpl"
}

func (mc *MainController) Delete() {
	key := KEY
	name := mc.GetSession("user")
	if name != nil {
		key = name.(string) + ":" + KEY
	} else {
		mc.Ctx.ResponseWriter.WriteHeader(403)
		mc.Data["json"] = map[string]string{"msg": "Does not logged in."}
		mc.ServeJSON()
		return
	}
	cli := NewClient()
	defer cli.Close()
	err := cli.Del(key).Err()
	if err != nil {
		panic(err)
		return
	}
	err = cli.Del(name.(string)).Err()
	if err != nil {
		panic(err)
		return
	}
	mc.Data["isLogin"] = false
	mc.Data["name"] = nil
	mc.DelSession("user")
	mc.TplName = "success.tpl"
}

func (mc *MainController) Clear() {
	key := KEY
	name := mc.GetSession("user")
	if name != nil {
		key = name.(string) + ":" + KEY
	}
	cli := NewClient()
	defer cli.Close()
	ret, err := cli.Del(key).Result()
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
	cli := NewClient()
	defer cli.Close()
	size, err := cli.LPush(key, m.Msg).Result()
	if err != nil {
		panic(err)
	}
	logs.Debug(size)
	mc.Redirect("/", 302)
}

func (mc *MainController) ListImage() {
	key := KEY
	name := mc.GetSession("user")
	if name != nil {
		key = name.(string) + ":" + KEY
		mc.Data["isLogin"] = true
		mc.Data["name"] = name
	}
	logs.Debug(name)
	cli := NewClient()
	defer cli.Close()
	images, err := cli.LRange(key, 0, -1).Result()
	if err != nil {
		panic(err)
	}
	logs.Debug(images)
	mc.Layout = "layout.tpl"
	mc.TplName = "board.tpl"
	mc.LayoutSections = make(map[string]string)
	mc.LayoutSections["Header"] = "header.tpl"
	mc.LayoutSections["Scripts"] = "scripts.tpl"
	mc.Data["images"] = images
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
	web.Router("/signup", new(MainController), "get:SignUp")
	web.Router("/create", new(MainController), "post:Create")
	web.Router("/signoff", new(MainController), "get:SignOff")
	web.Router("/delete", new(MainController), "post:Delete")
	web.Router("/image", new(MainController), "*:ListImage")
	web.AddFuncMap("rep", replace)
	web.Run()
}
