package admin

import (
	"amidesk-api-server/config"
	"amidesk-api-server/helper/rustdesk"

	"github.com/kataras/iris/v12/mvc"
)

type BootstrapController struct {
	basicController
}

func (c *BootstrapController) GetBootstrap() mvc.Result {
	return c.Success(rustdesk.BootstrapStatusFor(config.GetServerConfig().RustdeskBootstrap), "ok")
}