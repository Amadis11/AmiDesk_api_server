package middleware

import (
	"fmt"
	"amidesk-api-server/config"

	"github.com/kataras/iris/v12"
)

func RequestLogger() iris.Handler {
	return func(context iris.Context) {
		if config.GetServerConfig().DebugMode && context.Request().RequestURI != "/api/heartbeat" {
			requestInfo := fmt.Sprintf("▶ %s:%s", context.Method(), context.Request().RequestURI)
			body, _ := context.GetBody()
			context.Application().Logger().Info(requestInfo)
			for header, value := range context.Request().Header {
				fmt.Println(header+":", value)
			}
			fmt.Println(string(body))
		}

		context.Next()
	}
}
