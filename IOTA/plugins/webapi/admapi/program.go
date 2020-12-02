package admapi

import (
	"fmt"
	"net/http"

	"wasp/client"
	"wasp/packages/hashing"
	"wasp/packages/registry"
	"wasp/plugins/webapi/httperrors"

	"github.com/labstack/echo"
)

func addProgramEndpoints(adm *echo.Group) {
	adm.POST("/"+client.PutProgramRoute, handlePutProgram)
	adm.GET("/"+client.GetProgramMetadataRoute(":hash"), handleGetProgramMetadata)
}

func handlePutProgram(c echo.Context) error {
	var req client.PutProgramRequest
	var err error

	if err := c.Bind(&req); err != nil {
		return httperrors.BadRequest("Invalid request body")
	}

	if req.VMType == "" {
		return httperrors.BadRequest("vm_type is required")
	}
	if req.Description == "" {
		return httperrors.BadRequest("description is required")
	}
	if req.Code == nil || len(req.Code) == 0 {
		return httperrors.BadRequest("code is required (base64-encoded binary data)")
	}

	progHash, err := registry.SaveProgramCode(req.Code)
	if err != nil {
		return err
	}

	md := &registry.ProgramMetadata{
		ProgramHash: progHash,
		VMType:      req.VMType,
		Description: req.Description,
	}

	// TODO it is always overwritten!
	if err = md.Save(); err != nil {
		return err
	}

	log.Infof("Program metadata record has been saved. Program hash: %s, description: %s",
		md.ProgramHash.String(), md.Description)
	return c.JSON(http.StatusCreated, &client.PutProgramResponse{ProgramHash: &progHash})
}

func handleGetProgramMetadata(c echo.Context) error {
	progHash, err := hashing.HashValueFromBase58(c.Param("hash"))
	if err != nil {
		return httperrors.BadRequest("Invalid request body")
	}

	md, err := registry.GetProgramMetadata(&progHash)
	if err != nil {
		return err
	}
	if md == nil {
		return httperrors.NotFound(fmt.Sprintf("Program not found: %v", progHash.String()))
	}

	return c.JSON(http.StatusOK, &client.ProgramMetadata{
		VMType:      md.VMType,
		Description: md.Description,
	})
}
