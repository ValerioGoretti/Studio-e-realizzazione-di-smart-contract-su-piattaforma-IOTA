package webapi

import (
	"net"

	"wasp/plugins/webapi/admapi"
	"wasp/plugins/webapi/info"
	"wasp/plugins/webapi/request"
	"wasp/plugins/webapi/state"
)

func addEndpoints(adminWhitelist []net.IP) {
	info.AddEndpoints(Server)
	request.AddEndpoints(Server)
	state.AddEndpoints(Server)
	admapi.AddEndpoints(Server, adminWhitelist)
	log.Infof("added web api endpoints")
}
