package main

import (
	"context"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"os"
	"strconv"
	"strings"
	"time"
)

const KEY = "memo"
const IMAGE_KEY = "image"
const FILE_KEY = "file"

type Memo struct {
	Msg string `form:"msg"`
}

type Image struct {
	Base64Img string `form:"image"`
}

type File struct {
	Base64File string `form:"base64str"`
	Filename   string `form:"filename"`
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

type ErrorController struct {
	web.Controller
}

type RedisClient interface {
	LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd
	LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Close() error
}

var RedisClientFactory = func() RedisClient {
	return NewClient()
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

func (mc *MainController) ListText() {
	name := mc.GetSession("user")
	if name == nil {
		name = "anonymous"
	} else {
		mc.Data["isLogin"] = true
		mc.Data["name"] = name
	}
	key := name.(string) + ":" + KEY
	logs.Debug(name)
	cli := RedisClientFactory()
	defer cli.Close()
	ctx := mc.Ctx.Request.Context()
	memos, err := cli.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		panic(err)
	}
	logs.Debug(memos)
	mc.Layout = "meta/layout.tpl"
	mc.TplName = "text.tpl"
	mc.LayoutSections = make(map[string]string)
	mc.LayoutSections["Header"] = "meta/header.tpl"
	mc.LayoutSections["Scripts"] = "meta/scripts.tpl"
	mc.Data["memos"] = memos
}

func (mc *MainController) SaveText() {
	name := mc.GetSession("user")
	if name == nil {
		name = "anonymous"
	}
	key := name.(string) + ":" + KEY
	m := Memo{}
	if err := mc.ParseForm(&m); err != nil {
		panic(err)
	}
	cli := RedisClientFactory()
	defer cli.Close()
	ctx := mc.Ctx.Request.Context()
	texts, _ := cli.LRange(ctx, key, 0, -1).Result()
	if len(texts) > 9 {
		mc.Ctx.ResponseWriter.WriteHeader(403)
		mc.Data["json"] = map[string]string{"msg": "Exceeds the number of texts that can be uploaded."}
		mc.ServeJSON()
		return
	}
	size, err := cli.LPush(ctx, key, m.Msg).Result()
	if err != nil {
		panic(err)
	}
	logs.Debug(size)
	mc.Redirect("/", 302)
}

func (mc *MainController) ClearText() {
	name := mc.GetSession("user")
	if name == nil {
		name = "anonymous"
	}
	key := name.(string) + ":" + KEY
	cli := RedisClientFactory()
	defer cli.Close()
	ctx := mc.Ctx.Request.Context()
	ret, err := cli.Del(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	logs.Debug(ret)
	mc.Redirect("/", 302)
}

func (mc *MainController) ListImage() {
	name := mc.GetSession("user")
	if name == nil {
		name = "anonymous"
	}
	key := name.(string) + ":" + IMAGE_KEY
	logs.Debug(name)
	cli := RedisClientFactory()
	defer cli.Close()
	ctx := mc.Ctx.Request.Context()
	images, err := cli.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		panic(err)
	}
	logs.Debug(len(images))
	mc.Layout = "meta/layout.tpl"
	mc.TplName = "image.tpl"
	mc.LayoutSections = make(map[string]string)
	mc.LayoutSections["Header"] = "meta/header.tpl"
	mc.LayoutSections["Scripts"] = "meta/scripts.tpl"
	mc.Data["images"] = images
}

func (mc *MainController) SaveImage() {
	name := mc.GetSession("user")
	if name == nil {
		name = "anonymous"
	}
	key := name.(string) + ":" + IMAGE_KEY
	img := Image{}
	if err := mc.ParseForm(&img); err != nil {
		panic(err)
	}
	cli := RedisClientFactory()
	defer cli.Close()
	ctx := mc.Ctx.Request.Context()
	images, _ := cli.LRange(ctx, key, 0, -1).Result()
	if len(images) > 0 {
		mc.Ctx.ResponseWriter.WriteHeader(403)
		mc.Data["json"] = map[string]string{"msg": "Exceeds the number of images that can be uploaded."}
		mc.ServeJSON()
		return
	}
	replaced_img := replace(img.Base64Img, " ", "+")
	size, err := cli.LPush(ctx, key, replaced_img).Result()
	if err != nil {
		panic(err)
	}
	logs.Debug(size)
	mc.Redirect("/image", 302)
}

func (mc *MainController) ClearImage() {
	name := mc.GetSession("user")
	if name == nil {
		name = "anonymous"
	}
	key := name.(string) + ":" + IMAGE_KEY
	cli := RedisClientFactory()
	defer cli.Close()
	ctx := mc.Ctx.Request.Context()
	ret, err := cli.Del(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	logs.Debug(ret)
	mc.Redirect("/image", 302)
}

func (mc *MainController) ListFile() {
	name := mc.GetSession("user")
	if name == nil {
		name = "anonymous"
	} else {
		mc.Data["isLogin"] = true
		mc.Data["name"] = name
	}
	key := name.(string) + ":" + FILE_KEY
	logs.Debug(name)
	cli := RedisClientFactory()
	defer cli.Close()
	ctx := mc.Ctx.Request.Context()
	files, err := cli.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		panic(err)
	}
	logs.Debug(len(files))
	mc.Layout = "meta/layout.tpl"
	mc.TplName = "file.tpl"
	mc.LayoutSections = make(map[string]string)
	mc.LayoutSections["Header"] = "meta/header.tpl"
	mc.LayoutSections["Scripts"] = "meta/scripts.tpl"
	mc.Data["files"] = files
}

func (mc *MainController) SaveFile() {
	name := mc.GetSession("user")
	if name == nil {
		name = "anonymous"
	}
	key := name.(string) + ":" + FILE_KEY
	file := File{}
	if err := mc.ParseForm(&file); err != nil {
		panic(err)
	}
	cli := RedisClientFactory()
	defer cli.Close()
	ctx := mc.Ctx.Request.Context()
	files, _ := cli.LRange(ctx, key, 0, -1).Result()
	if len(files) > 0 {
		mc.Ctx.ResponseWriter.WriteHeader(403)
		mc.Data["json"] = map[string]string{"msg": "Exceeds the number of files that can be uploaded."}
		mc.ServeJSON()
		return
	}
	replaced_file := replace(file.Base64File, " ", "+")
	size, err := cli.LPush(ctx, key, replaced_file+"^_^"+file.Filename).Result()
	if err != nil {
		panic(err)
	}
	logs.Debug(size)
	mc.Redirect("/file", 302)
}

func (mc *MainController) ClearFile() {
	name := mc.GetSession("user")
	if name == nil {
		name = "anonymous"
	}
	key := name.(string) + ":" + FILE_KEY
	cli := RedisClientFactory()
	defer cli.Close()
	ctx := mc.Ctx.Request.Context()
	ret, err := cli.Del(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	logs.Debug(ret)
	mc.Redirect("/file", 302)
}

func (mc *MainController) Login() {
	mc.Layout = "meta/layout.tpl"
	mc.TplName = "account/login.tpl"
	mc.LayoutSections = make(map[string]string)
	mc.LayoutSections["Header"] = "meta/header.tpl"
	mc.LayoutSections["Scripts"] = "meta/scripts.tpl"
}

func (mc *MainController) SignUp() {
	mc.Layout = "meta/layout.tpl"
	mc.TplName = "account/signup.tpl"
	mc.LayoutSections = make(map[string]string)
	mc.LayoutSections["Header"] = "meta/header.tpl"
	mc.LayoutSections["Scripts"] = "meta/scripts.tpl"
}

func (mc *MainController) SignOff() {
	name := mc.GetSession("user")
	if name == nil {
		mc.Redirect("/", 302)
	}
	mc.Data["name"] = name
	mc.Layout = "meta/layout.tpl"
	mc.TplName = "account/signoff.tpl"
	mc.LayoutSections = make(map[string]string)
	mc.LayoutSections["Header"] = "meta/header.tpl"
	mc.LayoutSections["Scripts"] = "meta/scripts.tpl"
}

func (mc *MainController) Logout() {
	mc.Data["isLogin"] = false
	mc.Data["name"] = nil
	mc.DelSession("user")
	mc.TplName = "account/success.tpl"
}

func (mc *MainController) Check() {
	u := User{}
	if err := mc.ParseForm(&u); err != nil {
		panic(err)
	}
	cli := RedisClientFactory()
	defer cli.Close()
	ctx := mc.Ctx.Request.Context()
	pass, err := cli.Get(ctx, u.Name).Result()
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
	mc.TplName = "account/success.tpl"
}

func (mc *MainController) Create() {
	newu := NewUser{}
	if err := mc.ParseForm(&newu); err != nil {
		panic(err)
		return
	}
	cli := RedisClientFactory()
	defer cli.Close()
	ctx := mc.Ctx.Request.Context()
	err := cli.Get(ctx, newu.Name).Err()
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
	err = cli.Set(ctx, newu.Name, hash, 0).Err()
	if err != nil {
		panic(err)
		return
	}
	mc.SetSession("user", newu.Name)
	mc.TplName = "account/success.tpl"
}

func (mc *MainController) Delete() {
	name := mc.GetSession("user")
	if name == nil {
		mc.Ctx.ResponseWriter.WriteHeader(403)
		mc.Data["json"] = map[string]string{"msg": "Does not logged in."}
		mc.ServeJSON()
		return
	}
	key := name.(string) + ":" + KEY
	key_file := name.(string) + ":" + FILE_KEY
	key_image := name.(string) + ":" + IMAGE_KEY
	cli := RedisClientFactory()
	defer cli.Close()
	ctx := mc.Ctx.Request.Context()
	err := cli.Del(ctx, key).Err()
	if err != nil {
		panic(err)
		return
	}
	err = cli.Del(ctx, key_file).Err()
	if err != nil {
		panic(err)
		return
	}
	err = cli.Del(ctx, key_image).Err()
	if err != nil {
		panic(err)
		return
	}
	err = cli.Del(ctx, name.(string)).Err()
	if err != nil {
		panic(err)
		return
	}
	mc.Data["isLogin"] = false
	mc.Data["name"] = nil
	mc.DelSession("user")
	mc.TplName = "account/success.tpl"
}

func (ec *ErrorController) Error404() {
	ec.Layout = "meta/layout.tpl"
	ec.TplName = "error/404.tpl"
	ec.LayoutSections = make(map[string]string)
	ec.LayoutSections["Header"] = "meta/header.tpl"
	ec.LayoutSections["Scripts"] = "meta/scripts.tpl"
}

func replace(str string, from string, to string) (out string) {
	return strings.Replace(str, from, to, -1)
}

func is_first(index int) (flag bool) {
	return index%4 == 0
}

func is_end(index int) (flag bool) {
	return (index+1)%4 == 0
}

func get_file_name(file_info string) (filename string) {
	file_info_list := strings.Split(file_info, "^_^")
	return file_info_list[len(file_info_list)-1]
}

func get_content(file_info string) (content string) {
	file_info_list := strings.Split(file_info, "^_^")
	return "href=" + file_info_list[0]
}

func attr(s string) template.HTMLAttr {
	return template.HTMLAttr(s)
}

func safe(s string) template.HTML {
	return template.HTML(s)
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
	web.Router("/", new(MainController), "*:ListText")
	web.Router("/insert", new(MainController), "post:SaveText")
	web.Router("/clear", new(MainController), "post:ClearText")
	// for image
	web.Router("/image", new(MainController), "get:ListImage")
	web.Router("/save", new(MainController), "post:SaveImage")
	web.Router("/clear_img", new(MainController), "post:ClearImage")
	// for file
	web.Router("/file", new(MainController), "get:ListFile")
	web.Router("/upload", new(MainController), "post:SaveFile")
	web.Router("/clear_file", new(MainController), "post:ClearFile")
	// for user management
	web.Router("/login", new(MainController), "get:Login")
	web.Router("/check", new(MainController), "post:Check")
	web.Router("/logout", new(MainController), "post:Logout")
	web.Router("/signup", new(MainController), "get:SignUp")
	web.Router("/create", new(MainController), "post:Create")
	web.Router("/signoff", new(MainController), "get:SignOff")
	web.Router("/delete", new(MainController), "post:Delete")
	// for template functions
	web.AddFuncMap("rep", replace)
	web.AddFuncMap("is_first", is_first)
	web.AddFuncMap("is_end", is_end)
	web.AddFuncMap("get_file_name", get_file_name)
	web.AddFuncMap("get_content", get_content)
	web.AddFuncMap("safe", safe)
	web.AddFuncMap("attr", attr)
	// for error handling
	web.ErrorController(&ErrorController{})
	web.Run()
}
