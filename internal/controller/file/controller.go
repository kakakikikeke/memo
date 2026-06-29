package file

import (
	"context"

	"github.com/beego/beego/v2/core/logs"
	"github.com/kakakikikeke/memo/internal/controller"
	basectrl "github.com/kakakikikeke/memo/internal/controller"
	"github.com/kakakikikeke/memo/internal/database"
	"github.com/kakakikikeke/memo/internal/repository"
	"github.com/kakakikikeke/memo/internal/service"
)

type FileForm struct {
	Base64File string `form:"base64str"`
	Filename   string `form:"filename"`
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
	c.LogAccess("ListFile")
	username := c.CurrentUsername()
	ctx := c.Ctx.Request.Context()
	files, err := c.getMemoService().ListFile(ctx, username)
	if err != nil {
		logs.Error("ListFile failed: ", err)
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte("Internal Server Error"))
		return
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
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.Data["json"] = map[string]string{"msg": "Invalid request."}
		c.ServeJSON()
		return
	}
	ctx := c.Ctx.Request.Context()
	replacedFile := basectrl.Replace(file.Base64File, " ", "+")
	if !basectrl.IsValidFileDataURL(replacedFile) {
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.Data["json"] = map[string]string{"msg": "Invalid file data."}
		c.ServeJSON()
		return
	}
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
		logs.Error("ClearFile failed: ", err)
		c.Ctx.ResponseWriter.WriteHeader(500)
		c.Data["json"] = map[string]string{"msg": "Internal server error."}
		c.ServeJSON()
		return
	}
	c.Redirect("/file", 302)
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
