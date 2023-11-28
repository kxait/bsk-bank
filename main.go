package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"bsk-bank/handlers"
	"bsk-bank/lib"
	"bsk-bank/views"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/libsql/libsql-client-go/libsql"
)

func main() {
	godotenv.Load()
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		fmt.Fprintf(os.Stderr, "DATABASE_URL env var must be set")
		os.Exit(1)
	}

	db, err := sql.Open("libsql", databaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not connect to DB %s (%s)", databaseUrl, err)
		os.Exit(1)
	}

	// migrations etc
	lib.SetupDatabase(db)

	application := gin.Default()
	application.HTMLRender = &lib.TemplRender{}

	application.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "", views.Hello("jon"))
	})

	application.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "", views.Login())
	})

	application.POST("/login", handlers.LoginHandler(db))

	application.Run(":8080")
}
