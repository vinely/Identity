package identity

import (
	"github.com/gin-gonic/gin"
	"github.com/vinely/Identity/api"
)

var (
	version  = "0.0.1"
	basePath = "/api/v1"
)

// RunService - run server
func RunService() error {
	e := gin.Default()
	apiv1 := e.Group(basePath)
	api.RegisterIdentityAPI(apiv1)
	api.RegisterClaimAPI(apiv1)
	return e.Run(":8080")
}
