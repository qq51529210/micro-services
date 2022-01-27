package codes

import (
	"github.com/qq51529210/micro-services/auth/service"
	"github.com/qq51529210/web/router"
)

func post(ctx *router.Context) {
	_type := ctx.Request.FormValue("type")
	switch _type {
	case "phone":
	default:
		service.FormValueError(ctx, "type", _type)
	}
}