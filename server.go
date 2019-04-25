package Identity

import (
	"github.com/gin-gonic/gin"
)

var (
	version  = "0.0.1"
	basePath = "/api/v1"
)

func idget(c *gin.Context) {

}

// RunService - run server
func RunService() error {
	e := gin.Default()
	e.GET(basePath+"did/id/:id", idget)

	return e.Run(":8080")
}
