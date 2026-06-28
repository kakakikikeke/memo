package file

import (
	"context"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	basectrl "github.com/kakakikikeke/memo/internal/controller"
	"github.com/kakakikikeke/memo/internal/repository"
	"github.com/kakakikikeke/memo/internal/service"
)

type FileForm struct {
	Base64File string `form:"base64str"`
	Filename   string `form:"filename"`
}

type Controller struct {
	web.Controller
	memoService *service.MemoService
}

func NewController(svc *service.MemoService) *Controller {
	return &Controller{memoService: svc}
}

func (c *Controller) getMemoService() *service.MemoService {
	if c.memoService == nil {
		c.memoService = service.NewMemoService(repository.NewRedisRepository(basectrl.NewClient()))
	}
	return c.memoService
}

func (c *Controller) List() {
	c.LogAccess("ListFile")
	username := c.CurrentUsername()
	ctx := c.Ctx.Request.Context()
	files, err := c.getMemoService().ListFile(ctx, username)
	if err != nil {
		panic(err)
	}
	logs.Debug(len(files))
	c.RenderLayout("file.tpl")
	c.Data["files"] = files
	c.SetLoginContext(c.GetSession("user"))
}

func (c *Controller) Save() {
	c.LogAccess("SaveFile")
	username := c.CurrentUsername()
	file := FileForm{}
	if err := c.ParseForm(&file); err != nil {
		panic(err)
	}
	ctx := c.Ctx.Request.Context()
	replacedFile := basectrl.Replace(file.Base64File, " ", "+")
	value := replacedFile + "^_^" + file.Filename
	if err := c.getMemoService().SaveFileWithLimit(ctx, username, value); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(403)
		c.Data["json"] = map[string]string{"msg": err.Error()}
		c.ServeJSON()
		return
	}
	c.Redirect("/file", 302)
}

func (c *Controller) Clear() {
	c.LogAccess("ClearFile")
	username := c.CurrentUsername()
	ctx := c.Ctx.Request.Context()
	if err := c.getMemoService().ClearFile(ctx, username); err != nil {
		panic(err)
	}
	c.Redirect("/file", 302)
}

func (c *Controller) CurrentUsername() string {
	if name := c.GetSession("user"); name != nil {
		if username, ok := name.(string); ok {
			return username
		}
	}
	return "anonymous"
}

func (c *Controller) SetLoginContext(name interface{}) {
	if name == nil {
		c.Data["isLogin"] = false
		c.Data["name"] = nil
		return
	}
	c.Data["isLogin"] = true
	c.Data["name"] = name
}

func (c *Controller) RenderLayout(tplName string) {
	c.Layout = "meta/layout.tpl"
	c.TplName = tplName
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["Header"] = "meta/header.tpl"
	c.LayoutSections["Scripts"] = "meta/scripts.tpl"
}

func (c *Controller) LogAccess(action string) {
	name := c.CurrentUsername()
	ip := c.Ctx.Input.IP()
	logs.Info(action, " user=", name, " ip=", ip)
}

func (c *Controller) ListFile(ctx context.Context, username string) ([]string, error) {
	return c.getMemoService().ListFile(ctx, username)
}

func (c *Controller) SaveFileWithLimit(ctx context.Context, username, value string) error {
	return c.getMemoService().SaveFileWithLimit(ctx, username, value)
}

func (c *Controller) ClearFile(ctx context.Context, username string) error {
	return c.getMemoService().ClearFile(ctx, username)
}
