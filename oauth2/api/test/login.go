package test

import (
	"fmt"
	"html/template"
	"net/url"

	"github.com/gin-gonic/gin"
)

var (
	tp *template.Template
)

func init() {
	tp, _ = template.New("test").Parse(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>测试 oauth2 登录</title>
</head>
<body>
<a href="{{.}}">oauth2登录</a>
</body>
</html>`)
}

func login(ctx *gin.Context) {
	query := make(url.Values)
	query.Set("response_type", "code")
	query.Set("client_id", app)
	query.Set("scope", "readwrite")
	query.Set("state", state)
	query.Set("redirect_uri", fmt.Sprintf("%s/oauth2", host))
	redirectURL := fmt.Sprintf("%s/oauth2/authorize?%s", oauth2Host, query.Encode())
	tp.Execute(ctx.Writer, redirectURL)
}
