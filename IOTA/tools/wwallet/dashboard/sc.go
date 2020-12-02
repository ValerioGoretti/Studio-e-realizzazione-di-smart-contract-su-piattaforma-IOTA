package dashboard

import (
	"wasp/tools/wwallet/sc"

	"github.com/labstack/echo"
)

type SCDashboard interface {
	Config() *sc.Config
	AddEndpoints(e *echo.Echo)
	AddTemplates(r Renderer)
}
