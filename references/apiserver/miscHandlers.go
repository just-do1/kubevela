package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/oam-dev/kubevela/references/apiserver/apis"
	"github.com/oam-dev/kubevela/version"
)

// GetVersion will return version for dashboard
func (s *APIServer) GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, apis.Response{
		Code: http.StatusOK,
		Data: map[string]string{"version": version.VelaVersion},
	})
}
