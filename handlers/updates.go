package handlers

import (
	"bsk-bank/views"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDashboardUpdates(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "", views.Updates())
	}
}
