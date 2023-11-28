package lib

import (
	"bsk-bank/views"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorPage(ctx *gin.Context, errorMessage string) {
	ctx.HTML(http.StatusBadRequest, "", views.Error(errorMessage))
}
