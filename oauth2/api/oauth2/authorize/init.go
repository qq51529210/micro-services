package authorize

import (
	"auth/api/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init(g gin.IRouter) {
	g = g.Group("/authorize", middleware.CheckSession)
	//
	g.GET("", get)
	g.POST("", post)
}

type reqQuery struct {
	// 指定用于授权流程的响应类型，常见的值包括
	// code ，用于授权码授权流程
	// token ，用于隐式授权流程
	ResponseType string `form:"response_type" binding:"oneof=code token"`
	// 表示客户端应用程序的唯一标识符，由授权服务器分配给客户端
	ClientID string `form:"client_id" binding:"required,max=40"`
	// 表示客户端请求的访问权限范围
	Scope string `form:"scope" binding:"oneof=read write readwrite"`
	// 用于在授权请求和授权响应之间传递状态，以防止 CSRF 攻击
	State string `form:"state"`
	// 指定授权服务器将用户重定向到的客户端应用程序的回调 URL
	RedirectURI string `form:"redirect_uri" binding:"required,uri"`
}