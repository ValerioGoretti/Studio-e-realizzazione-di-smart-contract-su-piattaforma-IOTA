package admapi

import (
	"net/http"

	"wasp/client"
	"wasp/plugins/gracefulshutdown"

	"github.com/labstack/echo"
)

func addShutdownEndpoint(adm *echo.Group) {
	adm.GET("/"+client.ShutdownRoute, handleShutdown)
}

func handleShutdown(c echo.Context) error {
	log.Info("Received a shutdown request from WebAPI.")
	gracefulshutdown.Shutdown()
	return c.String(http.StatusOK, "Shutting down...")
}
