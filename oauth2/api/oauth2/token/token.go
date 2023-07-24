package token

import (
	"net/http"
	"net/url"
	"oauth2/api/internal"
	"oauth2/api/internal/middleware"
	"oauth2/db"

	"github.com/gin-gonic/gin"
	"github.com/qq51529210/util"
)

// 模式
const (
	GrantTypeAuthorizationCode = "authorization_code"
	GrantTypePassword          = "password"
	GrantTypeClientCredentials = "client_credentials"
	GrantTypeImplicit          = "implicit"
	GrantTypeRefreshToken      = "refresh_token"
)

// token 处理获取访问令牌
func token(ctx *gin.Context) {
	switch ctx.Query("grant_type") {
	case GrantTypeAuthorizationCode:
		tokenAuthorizationCode(ctx)
	case GrantTypeImplicit:
	case GrantTypePassword:
		tokenPassword(ctx)
	case GrantTypeClientCredentials:
	case GrantTypeRefreshToken:
	}
}

func onOK(ctx *gin.Context, token *db.AccessToken) {
	// 重定向
	redirectURI := ctx.Query(middleware.QueryRedirectURI)
	if redirectURI != "" {
		// 重定向地址
		_u, err := url.Parse(redirectURI)
		if err != nil {
			internal.Submit400(ctx, err.Error())
			return
		}
		_u.RawQuery = util.HTTPQuery(token, _u.Query()).Encode()
		// 跳转
		ctx.Redirect(http.StatusSeeOther, _u.String())
		return
	}
	// 没有重定向，返回 JSON
	ctx.JSON(http.StatusOK, token)
}
