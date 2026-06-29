package text

import (
	"context"
	"errors"

	"github.com/beego/beego/v2/core/logs"
	"github.com/kakakikikeke/memo/internal/controller"
	"github.com/kakakikikeke/memo/internal/database"
	"github.com/kakakikikeke/memo/internal/repository"
	"github.com/kakakikikeke/memo/internal/service"
)

type MemoForm struct {
	Msg string `form:"msg"`
}

type Controller struct {
	controller.BaseController
	memoService *service.MemoService
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
		logs.Error("ListText failed: ", err)
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte("Internal Server Error"))
		return
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
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.Data["json"] = map[string]string{"msg": "Invalid request."}
		c.ServeJSON()
		return
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
		logs.Error("ClearText failed: ", err)
		c.Ctx.ResponseWriter.WriteHeader(500)
		c.Data["json"] = map[string]string{"msg": "Internal server error."}
		c.ServeJSON()
		return
	}
	c.Redirect("/", 302)
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
