package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary		Ping-Pong API
// @Description	Check if server is up and running
// @Tags			server
// @Accept			json
// @Produce		json
// @Success		200	{string}	pong
// @Router			/v1/ping [get]
func (s Server) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
}
