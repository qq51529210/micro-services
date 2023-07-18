package middleware

import (
	"auth/api/internal"
	"auth/db"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

const (
	// CookieName 表示 cookie 的名称
	CookieName = "sid"
	// SessionContextKey 表示 session 上下文数据的 key
	SessionContextKey = "sck"
	// 登录 url
	loginURL = "/login"
	// QueryRedirectURL 重定向，查询参数名称
	QueryRedirectURL = "redirect_url"
)

// CheckSession 使用 cookie 检查用户登录
func CheckSession(ctx *gin.Context) {
	// 提取 cookie
	sid, err := ctx.Cookie(CookieName)
	if err == http.ErrNoCookie {
		redirectLogin(ctx)
		return
	}
	// 查询 session
	sess, err := db.GetSession(sid)
	if err != nil {
		internal.DB500(ctx, err)
		return
	}
	// 没有
	if sess == nil {
		redirectLogin(ctx)
		return
	}
	// 设置上下文
	ctx.Set(SessionContextKey, sess)
	//
	ctx.Next()
}

// redirectLogin 重定向到 /login
func redirectLogin(ctx *gin.Context) {
	query := make(url.Values)
	query.Set(QueryRedirectURL, ctx.Request.URL.RawQuery)
	redirectURL := loginURL + "?" + query.Encode()
	ctx.Redirect(http.StatusFound, redirectURL)
}