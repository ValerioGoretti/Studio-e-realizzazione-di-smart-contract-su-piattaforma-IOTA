package info

import (
	"net/http"

	"wasp/client"
	"wasp/packages/parameters"
	"wasp/plugins/banner"
	"wasp/plugins/peering"

	"github.com/labstack/echo"
)

func AddEndpoints(server *echo.Echo) {
	server.GET("/"+client.InfoRoute, handleInfo)
}

func handleInfo(c echo.Context) error {
	return c.JSON(http.StatusOK, client.InfoResponse{
		Version:       banner.AppVersion,
		NetworkId:     peering.MyNetworkId(),
		PublisherPort: parameters.GetInt(parameters.NanomsgPublisherPort),
	})
}
