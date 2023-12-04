package handlers

import (
	"bsk-bank/lib"
	"bsk-bank/views"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DashboardAccountsHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := lib.GetUser(ctx, db)
		if err != nil {
			views.ErrorPage(ctx, "Must be logged in!")
			return
		}

		fzy := ctx.Query("query")
		var records *sql.Rows

		if fzy != "" {
			fzy = fmt.Sprintf("%%%s%%", fzy)
			records, err = db.Query("SELECT * FROM account a WHERE a.holder_name LIKE ?", fzy)
		} else {
			records, err = db.Query("SELECT * FROM account")
		}
		if err != nil {
			views.ErrorPage(ctx, "Nie udało się pobrać kont!"+err.Error())
			return
		}

		accounts, err := lib.ScanAccounts(records)
		if err != nil {
			views.ErrorPage(ctx, "Nie udało się pobrać kont!"+err.Error())
			return
		}

		balances, err := lib.GetAccountsBalances(db)
		if err != nil {
			views.ErrorPage(ctx, "Nie udało się pobrać kont! "+err.Error())
			return
		}

		ctx.HTML(http.StatusOK, "", views.Accounts(views.AccountsListViewModel{Balances: balances, Accounts: accounts}))
	}
}

func GetDashboardCreateAccountHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := lib.GetUser(ctx, db)
		if err != nil {
			views.ErrorPage(ctx, "Must be logged in!")
			return
		}

		ctx.HTML(http.StatusOK, "", views.CreateAccount())
	}
}

func PostDashboardCreateAccountHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := lib.GetUser(ctx, db)
		if err != nil {
			views.ErrorPage(ctx, "Must be logged in!")
			return
		}

		name := ctx.PostForm("name")
		if name == "" {
			views.ErrorPage(ctx, "Podaj imie i nazwisko")
			return
		}

		records, err := db.Query("INSERT INTO account (holder_name) VALUES (?) RETURNING *", name)
		if err != nil {
			views.ErrorPage(ctx, err.Error())
			return
		}

		accounts, err := lib.ScanAccounts(records)
		if err != nil {
			views.ErrorPage(ctx, err.Error())
			return
		}

		if len(accounts) != 1 {
			views.ErrorPage(ctx, "Nie udalo sie stworzyc konta!")
			return
		}

		ctx.Redirect(http.StatusFound, "/dashboard/accounts")
	}
}
