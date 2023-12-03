package main

import (
	"database/sql"
	"fmt"
	"os"

	"bsk-bank/handlers"
	"bsk-bank/lib"

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
	err = lib.SetupDatabase(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "uh oh %s", err)
		os.Exit(1)
	}

	application := gin.Default()
	application.HTMLRender = &lib.TemplRender{}

	application.GET("/", handlers.HelloHandler(db))

	application.GET("/login", handlers.GetLoginHandler(db))
	application.POST("/login", handlers.PostLoginHandler(db))

	application.GET("/logout", handlers.LogoutHandler(db))

	application.GET("/dashboard", handlers.DashboardHandler(db))
	application.GET("/dashboard/accounts", handlers.DashboardAccountsHandler(db))
	application.GET("/dashboard/accounts/create", handlers.GetDashboardCreateAccountHandler(db))
	application.POST("/dashboard/accounts/create", handlers.PostDashboardCreateAccountHandler(db))
	application.GET("/dashboard/accounts/:accountId", handlers.DashboardAccountDetailsHandler(db))
	application.POST("/dashboard/accounts/:accountId/transaction", handlers.PostDashboardAccountDetailsCreateTransactionHandler(db))
	application.POST("/dashboard/accounts/:accountId/deposit", handlers.PostDashboardAccountDetailsDepositHandler(db))
	application.POST("/dashboard/accounts/:accountId/withdraw", handlers.PostDashboardAccountDetailsWithdrawHandler(db))
	application.GET("/dashboard/transactions", handlers.TransactionsHandler(db))

	application.Run(":8080")
}
