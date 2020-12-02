package admapi

import (
	"fmt"
	"net/http"

	"wasp/client"
	"wasp/packages/registry"
	"wasp/plugins/committees"
	"wasp/plugins/webapi/httperrors"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
	"github.com/labstack/echo"
)

func addSCEndpoints(adm *echo.Group) {
	adm.POST("/"+client.ActivateSCRoute(":address"), handleActivateSC)
	adm.POST("/"+client.DeactivateSCRoute(":address"), handleDeactivateSC)
}

func handleActivateSC(c echo.Context) error {
	scAddress, err := address.FromBase58(c.Param("address"))
	if err != nil {
		return httperrors.BadRequest(fmt.Sprintf("Invalid SC address: %s", c.Param("address")))
	}

	bd, err := registry.ActivateBootupData(&scAddress)
	if err != nil {
		return err
	}

	log.Debugw("calling committees.ActivateCommittee", "addr", bd.Address.String())
	if err := committees.ActivateCommittee(bd); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func handleDeactivateSC(c echo.Context) error {
	scAddress, err := address.FromBase58(c.Param("address"))
	if err != nil {
		return httperrors.BadRequest(fmt.Sprintf("Invalid SC address: %s", c.Param("address")))
	}

	bd, err := registry.DeactivateBootupData(&scAddress)
	if err != nil {
		return err
	}

	err = committees.DeactivateCommittee(bd)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
