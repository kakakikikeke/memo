package text

import (
	"context"
	"errors"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	basectrl "github.com/kakakikikeke/memo/internal/controller"
	"github.com/kakakikikeke/memo/internal/database"
	"github.com/kakakikikeke/memo/internal/repository"
	"github.com/kakakikikeke/memo/internal/service"
)

type MemoForm struct {
	Msg string `form:"msg"`
}

type Controller struct {
	web.Controller
	baseController *basectrl.BaseController
	memoService    *service.MemoService
}

func NewController(svc *service.MemoService) *Controller {
	return &Controller{memoService: svc}
}

func (c *Controller) getMemoService() *service.MemoService {
	if c.memoService == nil {
		c.memoService = service.NewMemoService(repository.NewRedisRepository(database.NewClient()))
	}
	return c.memoService
}

func (c *Controller) List() {
	c.LogAccess("ListText")
	username := c.CurrentUsername()
	ctx := c.Ctx.Request.Context()
	memos, err := c.getMemoService().ListText(ctx, username)
	if err != nil {
		panic(err)
	}
	logs.Debug(memos)
	c.RenderLayout("text.tpl")
	c.Data["memos"] = memos
	c.SetLoginContext(c.GetSession("user"))
}

func (c *Controller) Save() {
	c.LogAccess("SaveText")
	username := c.CurrentUsername()
	m := MemoForm{}
	if err := c.ParseForm(&m); err != nil {
		panic(err)
	}
	ctx := c.Ctx.Request.Context()
	if err := c.getMemoService().SaveTextWithLimit(ctx, username, m.Msg); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(403)
		c.Data["json"] = map[string]string{"msg": err.Error()}
		c.ServeJSON()
		return
	}
	c.Redirect("/", 302)
}

func (c *Controller) Clear() {
	c.LogAccess("ClearText")
	username := c.CurrentUsername()
	ctx := c.Ctx.Request.Context()
	if err := c.getMemoService().ClearText(ctx, username); err != nil {
		panic(err)
	}
	c.Redirect("/", 302)
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

func (c *Controller) SaveWithLimit(ctx context.Context, username, msg string) error {
	return c.getMemoService().SaveTextWithLimit(ctx, username, msg)
}

func (c *Controller) ClearText(ctx context.Context, username string) error {
	return c.getMemoService().ClearText(ctx, username)
}

func (c *Controller) ListText(ctx context.Context, username string) ([]string, error) {
	return c.getMemoService().ListText(ctx, username)
}

func (c *Controller) Error(err error) {
	if errors.Is(err, service.ErrTextLimitExceeded) {
		c.Ctx.ResponseWriter.WriteHeader(403)
		c.Data["json"] = map[string]string{"msg": err.Error()}
		c.ServeJSON()
	}
}
