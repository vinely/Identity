package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vinely/Identity/ont"
)

func listClaimBySubject(c *gin.Context) {
	strPage := c.Query("page")
	page, err := strconv.ParseUint(strPage, 10, 32)
	if err != nil {
		page = 0
	}
	id := c.Param("id")

	data, err := ont.ClaimList(uint(page), func(k, v []byte) (*ont.ClaimStatus, error) {
		clm := &ont.ClaimStatus{}
		err := clm.Unmarshal(v)
		if err != nil {
			return nil, err
		}
		if id == "" {
			return clm, nil
		}
		hd, err := clm.Hdr()
		if id == hd.ID() {
			return clm, nil
		}
		return nil, errors.New("not matched")

	})
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"Identities": data})
}

func addClaim(c *gin.Context) {

}

// RegisterClaimAPI - register claim api to gin engine
func RegisterClaimAPI(r *gin.RouterGroup) {

	r.GET("/claim/subject/:id", listClaimBySubject)
	r.PUT("/claim", addClaim)
}
