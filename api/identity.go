package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vinely/Identity/ont"
	"github.com/vinely/dids"
	chain "github.com/vinely/ontchain"
)

func addIdentity(c *gin.Context) {
	id := c.Param("id")
	var (
		identity *ont.Identity
		mi       *chain.ManagedIdentity
		err      error
	)
	if id == "" {
		identity, mi, err = ont.NewIdentity()
	} else {
		did := dids.ID(id)
		identity, mi, err = ont.IdentityFromID(&did)
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ont.SaveManagedIdentity(mi)
	identity.SaveToDB()
	c.JSON(http.StatusOK, identity)
}

func listIdentity(c *gin.Context) {
	strPage := c.Query("page")
	page, err := strconv.ParseUint(strPage, 10, 32)
	if err != nil {
		page = 0
	}

	data, err := ont.IdentityList(uint(page), func(k, v []byte) (*ont.Identity, error) {
		id := &ont.Identity{}
		err := id.Unmarshal(v)
		if err != nil {
			return nil, err
		}
		return id, nil
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"Identities": data})
}

func readIdentity(c *gin.Context) {
	id := c.Param("id")
	identity, err := ont.GetIdentityFromDB(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, identity)
}

func getTotalIdentity(c *gin.Context) {
	c.JSON(http.StatusOK, ont.IdentityNumber())
}

func readAdminIdentity(c *gin.Context) {
	c.JSON(http.StatusOK, ont.Admin)
}

func updateIdentity(c *gin.Context) {
	i := &ont.Identity{}
	c.BindJSON(i)
	err := i.SaveToDB()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, i)
}

// RegisterIdentityAPI - register Identit api to gin engine
func RegisterIdentityAPI(r *gin.RouterGroup) {

	r.PUT("/identity", addIdentity)
	r.GET("/identity", listIdentity)
	r.GET("/identity/id/:id", readIdentity)
	r.GET("/identity/total", getTotalIdentity)
	r.GET("/identity/admin", readAdminIdentity)

	// need authentication
	r.POST("/identity/id/:id", updateIdentity)

}
