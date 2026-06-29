package controller

import (
	"encoding/base64"
	"html/template"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

const (
	KeyMemo  = "memo"
	KeyImage = "image"
	KeyFile  = "file"
)

type BaseController struct {
	web.Controller
}

type ErrorController struct {
	web.Controller
}

func (c *BaseController) CurrentUsername() string {
	if name := c.GetSession("user"); name != nil {
		if username, ok := name.(string); ok {
			return username
		}
	}
	return "anonymous"
}

func (c *BaseController) SetLoginContext(name interface{}) {
	if name == nil {
		c.Data["isLogin"] = false
		c.Data["name"] = nil
		return
	}
	c.Data["isLogin"] = true
	c.Data["name"] = name
}

func (c *BaseController) RenderLayout(tplName string) {
	c.Layout = "meta/layout.tpl"
	c.TplName = tplName
	c.Data["csrf_token"] = c.XSRFToken()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["Header"] = "meta/header.tpl"
	c.LayoutSections["Scripts"] = "meta/scripts.tpl"
}

func (c *BaseController) LogAccess(action string) {
	name := c.CurrentUsername()
	ip := c.Ctx.Input.IP()
	logs.Info(action, " user=", name, " ip=", ip)
}

func (c *ErrorController) Error404() {
	c.Layout = "meta/layout.tpl"
	c.TplName = "error/404.tpl"
	c.Data["csrf_token"] = c.XSRFToken()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["Header"] = "meta/header.tpl"
	c.LayoutSections["Scripts"] = "meta/scripts.tpl"
}

func Replace(str string, from string, to string) string {
	return strings.Replace(str, from, to, -1)
}

func IsFirst(index int) bool {
	return index%4 == 0
}

func IsEnd(index int) bool {
	return (index+1)%4 == 0
}

func GetFileName(fileInfo string) string {
	fileInfoList := strings.SplitN(fileInfo, "^_^", 2)
	if len(fileInfoList) < 2 {
		return ""
	}
	return fileInfoList[len(fileInfoList)-1]
}

func GetContent(fileInfo string) string {
	fileInfoList := strings.SplitN(fileInfo, "^_^", 2)
	if len(fileInfoList) < 2 {
		return ""
	}
	content := fileInfoList[0]
	if !IsValidFileDataURL(content) {
		return ""
	}
	return content
}

func IsValidFileDataURL(value string) bool {
	const prefix = "data:"
	const marker = ";base64,"

	if !strings.HasPrefix(value, prefix) {
		return false
	}
	idx := strings.Index(value, marker)
	if idx <= len(prefix) {
		return false
	}
	mime := strings.ToLower(value[len(prefix):idx])
	if !isAllowedFileMIME(mime) {
		return false
	}
	encoded := value[idx+len(marker):]
	if encoded == "" {
		return false
	}
	if _, err := base64.StdEncoding.DecodeString(encoded); err == nil {
		return true
	}
	_, err := base64.RawStdEncoding.DecodeString(encoded)
	return err == nil
}

func isAllowedFileMIME(mime string) bool {
	allowed := map[string]struct{}{
		"application/octet-stream": {},
		"application/pdf":          {},
		"application/zip":          {},
		"image/gif":                {},
		"image/jpeg":               {},
		"image/png":                {},
		"text/csv":                 {},
		"text/plain":               {},
	}
	_, ok := allowed[mime]
	return ok
}

func Safe(s string) template.HTML {
	return template.HTML(s)
}
