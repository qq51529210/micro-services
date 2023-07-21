package token

import (
	"net/http"
	"net/url"
	"oauth2/api/internal"
	"oauth2/db"

	"github.com/gin-gonic/gin"
	"github.com/qq51529210/util"
)

type authorizationCodeReq struct {
	// 在授权码模式中使用，表示从授权服务器获取的授权码
	Code string `form:"code" binding:"required"`
	// 表示客户端应用程序的唯一标识符，由授权服务器分配给客户端
	ClientID string `form:"client_id" binding:"required,max=40"`
	// 表示客户端应用程序的密钥，由授权服务器分配给客户端
	ClientSecret string `form:"client_secret" binding:"required,max=40"`
	// 重定向 URL
	RedirectURI string `form:"redirect_uri" binding:"uri"`
}

// authorizationCode 处理 grant_type=authorization_code
func authorizationCode(ctx *gin.Context) {
	// 参数
	var req authorizationCodeReq
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		internal.Submit400(ctx, err.Error())
		return
	}
	// 授权码
	code, err := db.GetAuthorizationCode(req.Code)
	if err != nil {
		internal.DB500(ctx, err)
		return
	}
	if code == nil || req.ClientID != code.ClientID {
		internal.Submit400(ctx, err.Error())
		return
	}
	// 应用
	client, err := db.GetClient(req.ClientID)
	if err != nil {
		internal.DB500(ctx, err)
		return
	}
	if client == nil ||
		*client.Enable == db.True ||
		*client.Secret == util.SHA1String(req.ClientSecret) {
		internal.Submit400(ctx, err.Error())
		return
	}
	// 令牌
	token := new(db.AccessToken)
	token.Type = *client.TokenType
	token.Scope = code.Scope
	token.ClientID = code.ClientID
	token.UserID = code.UserID
	err = db.PutAccessToken(token)
	if err != nil {
		internal.DB500(ctx, err)
		return
	}
	// 重定向
	if req.RedirectURI != "" {
		// 重定向地址
		_u, err := url.Parse(req.RedirectURI)
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