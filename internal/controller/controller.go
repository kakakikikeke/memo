package controller

import (
	"errors"
	"html/template"
	"os"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/kakakikikeke/memo/internal/repository"
	"github.com/kakakikikeke/memo/internal/service"
	"github.com/redis/go-redis/v9"
)

const (
	KeyMemo  = "memo"
	KeyImage = "image"
	KeyFile  = "file"
)

type MemoForm struct {
	Msg string `form:"msg"`
}

type ImageForm struct {
	Base64Img string `form:"image"`
}

type FileForm struct {
	Base64File string `form:"base64str"`
	Filename   string `form:"filename"`
}

type UserForm struct {
	Name string `form:"name"`
	Pass string `form:"pass"`
}

type NewUserForm struct {
	Name  string `form:"name"`
	Pass  string `form:"pass"`
	Pass2 string `form:"pass2"`
}

type MainController struct {
	web.Controller
	memoService *service.MemoService
}

type ErrorController struct {
	web.Controller
}

type RedisClient = repository.RedisClient

var RedisClientFactory = func() RedisClient {
	return NewClient()
}

func (mc *MainController) getMemoService() *service.MemoService {
	if mc.memoService == nil {
		mc.memoService = service.NewMemoService(repository.NewRedisRepository(RedisClientFactory()))
	}
	return mc.memoService
}

func (mc *MainController) SetMemoService(svc *service.MemoService) {
	mc.memoService = svc
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
	logAccessMain(mc, "ListText")
	name := mc.GetSession("user")
	if name == nil {
		name = "anonymous"
	} else {
		mc.Data["isLogin"] = true
		mc.Data["name"] = name
	}
	username := name.(string)
	ctx := mc.Ctx.Request.Context()
	memos, err := mc.getMemoService().ListText(ctx, username)
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
	logAccessMain(mc, "SaveText")
	name := mc.GetSession("user")
	if name == nil {
		name = "anonymous"
	}
	username := name.(string)
	m := MemoForm{}
	if err := mc.ParseForm(&m); err != nil {
		panic(err)
	}
	ctx := mc.Ctx.Request.Context()
	if err := mc.getMemoService().SaveTextWithLimit(ctx, username, m.Msg); err != nil {
		mc.Ctx.ResponseWriter.WriteHeader(403)
		mc.Data["json"] = map[string]string{"msg": err.Error()}
		mc.ServeJSON()
		return
	}
	mc.Redirect("/", 302)
}

func (mc *MainController) ClearText() {
	logAccessMain(mc, "ClearText")
	name := mc.GetSession("user")
	if name == nil {
		name = "anonymous"
	}
	username := name.(string)
	ctx := mc.Ctx.Request.Context()
	if err := mc.getMemoService().ClearText(ctx, username); err != nil {
		panic(err)
	}
	mc.Redirect("/", 302)
}

func (mc *MainController) ListImage() {
	name := mc.GetSession("user")
	if name == nil {
		name = "anonymous"
	}
	username := name.(string)
	ctx := mc.Ctx.Request.Context()
	images, err := mc.getMemoService().ListImage(ctx, username)
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

func IsFirst(index int) (flag bool) {
	return index%4 == 0
}

func IsEnd(index int) (flag bool) {
	return (index+1)%4 == 0
}

func GetFileName(fileInfo string) (filename string) {
	fileInfoList := strings.Split(fileInfo, "^_^")
	return fileInfoList[len(fileInfoList)-1]
}

func GetContent(fileInfo string) (content string) {
	fileInfoList := strings.Split(fileInfo, "^_^")
	return "href=" + fileInfoList[0]
}

func Attr(s string) template.HTMLAttr {
	return template.HTMLAttr(s)
}

func Safe(s string) template.HTML {
	return template.HTML(s)
}

func (mc *MainController) SaveImage() {
	name := mc.GetSession("user")
	if name == nil {
		name = "anonymous"
	}
	username := name.(string)
	img := ImageForm{}
	if err := mc.ParseForm(&img); err != nil {
		panic(err)
	}
	ctx := mc.Ctx.Request.Context()
	replacedImg := Replace(img.Base64Img, " ", "+")
	if err := mc.getMemoService().SaveImageWithLimit(ctx, username, replacedImg); err != nil {
		mc.Ctx.ResponseWriter.WriteHeader(403)
		mc.Data["json"] = map[string]string{"msg": err.Error()}
		mc.ServeJSON()
		return
	}
	mc.Redirect("/image", 302)
}

func (mc *MainController) ClearImage() {
	name := mc.GetSession("user")
	if name == nil {
		name = "anonymous"
	}
	username := name.(string)
	ctx := mc.Ctx.Request.Context()
	if err := mc.getMemoService().ClearImage(ctx, username); err != nil {
		panic(err)
	}
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
	username := name.(string)
	ctx := mc.Ctx.Request.Context()
	files, err := mc.getMemoService().ListFile(ctx, username)
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
	username := name.(string)
	file := FileForm{}
	if err := mc.ParseForm(&file); err != nil {
		panic(err)
	}
	ctx := mc.Ctx.Request.Context()
	replacedFile := Replace(file.Base64File, " ", "+")
	value := replacedFile + "^_^" + file.Filename
	if err := mc.getMemoService().SaveFileWithLimit(ctx, username, value); err != nil {
		mc.Ctx.ResponseWriter.WriteHeader(403)
		mc.Data["json"] = map[string]string{"msg": err.Error()}
		mc.ServeJSON()
		return
	}
	mc.Redirect("/file", 302)
}

func (mc *MainController) ClearFile() {
	name := mc.GetSession("user")
	if name == nil {
		name = "anonymous"
	}
	username := name.(string)
	ctx := mc.Ctx.Request.Context()
	if err := mc.getMemoService().ClearFile(ctx, username); err != nil {
		panic(err)
	}
	mc.Redirect("/file", 302)
}

func (mc *MainController) Login() {
	logAccessMain(mc, "Login")
	mc.Layout = "meta/layout.tpl"
	mc.TplName = "account/login.tpl"
	mc.LayoutSections = make(map[string]string)
	mc.LayoutSections["Header"] = "meta/header.tpl"
	mc.LayoutSections["Scripts"] = "meta/scripts.tpl"
}

func (mc *MainController) SignUp() {
	logAccessMain(mc, "SignUp")
	mc.Layout = "meta/layout.tpl"
	mc.TplName = "account/signup.tpl"
	mc.LayoutSections = make(map[string]string)
	mc.LayoutSections["Header"] = "meta/header.tpl"
	mc.LayoutSections["Scripts"] = "meta/scripts.tpl"
}

func (mc *MainController) SignOff() {
	logAccessMain(mc, "SignOff")
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
	logAccessMain(mc, "Logout")
	mc.Data["isLogin"] = false
	mc.Data["name"] = nil
	mc.DelSession("user")
	mc.TplName = "account/success.tpl"
}

func (mc *MainController) Check() {
	logAccessMain(mc, "Check")
	u := UserForm{}
	if err := mc.ParseForm(&u); err != nil {
		panic(err)
	}
	ctx := mc.Ctx.Request.Context()
	ok, err := mc.getMemoService().Authenticate(ctx, u.Name, u.Pass)
	if err != nil {
		mc.Ctx.ResponseWriter.WriteHeader(403)
		mc.Data["json"] = map[string]string{"msg": "User does not found."}
		mc.ServeJSON()
		return
	}
	if !ok {
		mc.Ctx.ResponseWriter.WriteHeader(403)
		mc.Data["json"] = map[string]string{"msg": "Login failed."}
		mc.ServeJSON()
		return
	}
	mc.SetSession("user", u.Name)
	mc.TplName = "account/success.tpl"
}

func (mc *MainController) Create() {
	logAccessMain(mc, "Create")
	newu := NewUserForm{}
	if err := mc.ParseForm(&newu); err != nil {
		panic(err)
	}
	ctx := mc.Ctx.Request.Context()
	if err := mc.getMemoService().CreateUser(ctx, newu.Name, newu.Pass); err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			mc.Ctx.ResponseWriter.WriteHeader(403)
			mc.Data["json"] = map[string]string{"msg": "Specified user already exists."}
			mc.ServeJSON()
			return
		}
		panic(err)
	}

	if newu.Pass != newu.Pass2 {
		mc.Ctx.ResponseWriter.WriteHeader(403)
		mc.Data["json"] = map[string]string{"msg": "Does not match password."}
		mc.ServeJSON()
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
	ctx := mc.Ctx.Request.Context()
	if err := mc.getMemoService().DeleteUser(ctx, name.(string)); err != nil {
		panic(err)
	}
	mc.Data["isLogin"] = false
	mc.Data["name"] = nil
	mc.DelSession("user")
	mc.TplName = "account/success.tpl"
}

func (ec *ErrorController) Error404() {
	logAccessError(ec, "Error404")
	ec.Layout = "meta/layout.tpl"
	ec.TplName = "error/404.tpl"
	ec.LayoutSections = make(map[string]string)
	ec.LayoutSections["Header"] = "meta/header.tpl"
	ec.LayoutSections["Scripts"] = "meta/scripts.tpl"
}

func Replace(str string, from string, to string) (out string) {
	return strings.Replace(str, from, to, -1)
}

func logAccessMain(mc *MainController, action string) {
	name := mc.GetSession("user")
	if name == nil {
		name = "anonymous"
	}
	ip := mc.Ctx.Input.IP()
	logs.Info(action, " user=", name, " ip=", ip)
}

func logAccessError(ec *ErrorController, action string) {
	ip := ec.Ctx.Input.IP()
	logs.Info(action, " ip=", ip)
}

func init() {
	_ = strconv.IntSize
}
