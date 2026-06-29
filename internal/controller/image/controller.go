package image

import (
	"context"

	"github.com/beego/beego/v2/core/logs"
	"github.com/kakakikikeke/memo/internal/controller"
	basectrl "github.com/kakakikikeke/memo/internal/controller"
	"github.com/kakakikikeke/memo/internal/database"
	"github.com/kakakikikeke/memo/internal/repository"
	"github.com/kakakikikeke/memo/internal/service"
)

type ImageForm struct {
	Base64Img string `form:"image"`
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
	c.LogAccess("ListImage")
	username := c.CurrentUsername()
	ctx := c.Ctx.Request.Context()
	images, err := c.getMemoService().ListImage(ctx, username)
	if err != nil {
		logs.Error("ListImage failed: ", err)
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte("Internal Server Error"))
		return
	}
	logs.Debug(len(images))
	c.RenderLayout("image.tpl")
	c.Data["images"] = images
	c.SetLoginContext(c.GetSession("user"))
}

func (c *Controller) Save() {
	c.LogAccess("SaveImage")
	username := c.CurrentUsername()
	img := ImageForm{}
	if err := c.ParseForm(&img); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.Data["json"] = map[string]string{"msg": "Invalid request."}
		c.ServeJSON()
		return
	}
	ctx := c.Ctx.Request.Context()
	replacedImg := basectrl.Replace(img.Base64Img, " ", "+")
	if err := c.getMemoService().SaveImageWithLimit(ctx, username, replacedImg); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(403)
		c.Data["json"] = map[string]string{"msg": err.Error()}
		c.ServeJSON()
		return
	}
	c.Redirect("/image", 302)
}

func (c *Controller) Clear() {
	c.LogAccess("ClearImage")
	username := c.CurrentUsername()
	ctx := c.Ctx.Request.Context()
	if err := c.getMemoService().ClearImage(ctx, username); err != nil {
		logs.Error("ClearImage failed: ", err)
		c.Ctx.ResponseWriter.WriteHeader(500)
		c.Data["json"] = map[string]string{"msg": "Internal server error."}
		c.ServeJSON()
		return
	}
	c.Redirect("/image", 302)
}

func (c *Controller) ListImage(ctx context.Context, username string) ([]string, error) {
	return c.getMemoService().ListImage(ctx, username)
}

func (c *Controller) SaveImageWithLimit(ctx context.Context, username, base64 string) error {
	return c.getMemoService().SaveImageWithLimit(ctx, username, base64)
}

func (c *Controller) ClearImage(ctx context.Context, username string) error {
	return c.getMemoService().ClearImage(ctx, username)
}
