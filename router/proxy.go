package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/types"
)

func addProxyRoutes(r *gin.RouterGroup) {
	r.GET("/redirects", handleGetRedirects)
	r.POST("/redirect", handleAddRedirect)
	r.DELETE("/redirect/:id", handleRemoveRedirect)
}

// handleGetRedirects handles the retrieval of all redirects.
func handleGetRedirects(c *gin.Context) {
	redirects := proxyService.GetRedirects()
	c.JSON(http.StatusOK, redirects)
}

type handleAddRedirectBody struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// handleAddRedirect handles the addition of a redirect.
// Errors can be:
//   - failed_to_parse_body: failed to parse the request body.
//   - failed_to_add_redirect: failed to add the redirect.
func handleAddRedirect(c *gin.Context) {
	var body handleAddRedirectBody
	err := c.BindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, types.APIError{
			Code:    "failed_to_parse_body",
			Message: fmt.Sprintf("failed to parse request body: %v", err),
		})
		return
	}

	redirect := types.ProxyRedirect{
		Source: body.Source,
		Target: body.Target,
	}

	err = proxyService.AddRedirect(redirect)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, types.APIError{
			Code:    "failed_to_add_redirect",
			Message: fmt.Sprintf("failed to add redirect: %v", err),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// handleRemoveRedirect handles the removal of a redirect.
// Errors can be:
//   - missing_redirect_uuid: missing redirect uuid.
//   - invalid_redirect_uuid: invalid redirect uuid.
//   - failed_to_remove_redirect: failed to remove the redirect.
func handleRemoveRedirect(c *gin.Context) {
	idString := c.Param("id")
	if idString == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, types.APIError{
			Code:    "missing_redirect_uuid",
			Message: "missing redirect uuid",
		})
		return
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, types.APIError{
			Code:    "invalid_redirect_uuid",
			Message: "invalid redirect uuid",
		})
		return
	}

	err = proxyService.RemoveRedirect(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, types.APIError{
			Code:    "failed_to_remove_redirect",
			Message: fmt.Sprintf("failed to remove redirect: %v", err),
		})
		return
	}

	c.Status(http.StatusNoContent)
}