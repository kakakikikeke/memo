package user

import (
	"context"
	"errors"

	"github.com/beego/beego/v2/core/logs"
	"github.com/kakakikikeke/memo/internal/controller"
	"github.com/kakakikikeke/memo/internal/database"
	"github.com/kakakikikeke/memo/internal/repository"
	"github.com/kakakikikeke/memo/internal/service"
)

type UserForm struct {
	Name string `form:"name"`
	Pass string `form:"pass"`
}

type NewUserForm struct {
	Name  string `form:"name"`
	Pass  string `form:"pass"`
	Pass2 string `form:"pass2"`
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

func (c *Controller) Login() {
	c.LogAccess("Login")
	c.RenderLayout("account/login.tpl")
}

func (c *Controller) SignUp() {
	c.LogAccess("SignUp")
	c.RenderLayout("account/signup.tpl")
}

func (c *Controller) SignOff() {
	c.LogAccess("SignOff")
	if name := c.GetSession("user"); name == nil {
		c.Redirect("/", 302)
		return
	}
	c.Data["name"] = c.GetSession("user")
	c.RenderLayout("account/signoff.tpl")
}

func (c *Controller) Logout() {
	c.LogAccess("Logout")
	c.Data["isLogin"] = false
	c.Data["name"] = nil
	c.DelSession("user")
	c.TplName = "account/success.tpl"
}

func (c *Controller) Check() {
	c.LogAccess("Check")
	u := UserForm{}
	if err := c.ParseForm(&u); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.Data["json"] = map[string]string{"msg": "Invalid request."}
		c.ServeJSON()
		return
	}
	ctx := c.Ctx.Request.Context()
	ok, err := c.getMemoService().Authenticate(ctx, u.Name, u.Pass)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(403)
		c.Data["json"] = map[string]string{"msg": "User does not found."}
		c.ServeJSON()
		return
	}
	if !ok {
		c.Ctx.ResponseWriter.WriteHeader(403)
		c.Data["json"] = map[string]string{"msg": "Login failed."}
		c.ServeJSON()
		return
	}
	c.SetSession("user", u.Name)
	c.TplName = "account/success.tpl"
}

func (c *Controller) Create() {
	c.LogAccess("Create")
	newu := NewUserForm{}
	if err := c.ParseForm(&newu); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.Data["json"] = map[string]string{"msg": "Invalid request."}
		c.ServeJSON()
		return
	}
	if newu.Pass != newu.Pass2 {
		c.Ctx.ResponseWriter.WriteHeader(403)
		c.Data["json"] = map[string]string{"msg": "Does not match password."}
		c.ServeJSON()
		return
	}
	ctx := c.Ctx.Request.Context()
	if err := c.getMemoService().CreateUser(ctx, newu.Name, newu.Pass); err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			c.Ctx.ResponseWriter.WriteHeader(403)
			c.Data["json"] = map[string]string{"msg": "Specified user already exists."}
			c.ServeJSON()
			return
		}
		logs.Error("CreateUser failed: ", err)
		c.Ctx.ResponseWriter.WriteHeader(500)
		c.Data["json"] = map[string]string{"msg": "Internal server error."}
		c.ServeJSON()
		return
	}
	c.SetSession("user", newu.Name)
	c.TplName = "account/success.tpl"
}

func (c *Controller) Delete() {
	name := c.GetSession("user")
	if name == nil {
		c.Ctx.ResponseWriter.WriteHeader(403)
		c.Data["json"] = map[string]string{"msg": "Does not logged in."}
		c.ServeJSON()
		return
	}
	ctx := c.Ctx.Request.Context()
	if err := c.getMemoService().DeleteUser(ctx, name.(string)); err != nil {
		logs.Error("DeleteUser failed: ", err)
		c.Ctx.ResponseWriter.WriteHeader(500)
		c.Data["json"] = map[string]string{"msg": "Internal server error."}
		c.ServeJSON()
		return
	}
	c.Data["isLogin"] = false
	c.Data["name"] = nil
	c.DelSession("user")
	c.TplName = "account/success.tpl"
}

func (c *Controller) Authenticate(ctx context.Context, username, password string) (bool, error) {
	return c.getMemoService().Authenticate(ctx, username, password)
}

func (c *Controller) CreateUser(ctx context.Context, username, password string) error {
	return c.getMemoService().CreateUser(ctx, username, password)
}

func (c *Controller) DeleteUser(ctx context.Context, username string) error {
	return c.getMemoService().DeleteUser(ctx, username)
}
