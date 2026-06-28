package model

import "fmt"

const (
	KeyMemo  = "memo"
	KeyImage = "image"
	KeyFile  = "file"
)

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

func UserKey(username, kind string) string {
	return fmt.Sprintf("%s:%s", username, kind)
}
