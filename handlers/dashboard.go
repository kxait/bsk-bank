package handlers

import (
	"bsk-bank/lib"
	"bsk-bank/views"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func DashboardHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := lib.GetUser(ctx, db)
		if err != nil {
			views.ErrorPage(ctx, "Must be logged in!")
			return
		}

		ctx.HTML(http.StatusOK, "", views.Dashboard())
	}
}

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

func DashboardAccountDetailsHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := lib.GetUser(ctx, db)
		if err != nil {
			views.ErrorPage(ctx, "Must be logged in!")
			return
		}

		accountId := ctx.Param("accountId")
		accountId = strings.ReplaceAll(accountId, "/", "")
		accountIdNum, err := strconv.ParseInt(accountId, 10, 64)

		if err != nil {
			views.ErrorPage(ctx, "ups"+err.Error())
			return
		}

		records, err := db.Query("SELECT * FROM account WHERE id = ?", accountIdNum)
		if err != nil {
			views.ErrorPage(ctx, "ups"+err.Error())
			return
		}

		accounts, err := lib.ScanAccounts(records)
		if err != nil {
			views.ErrorPage(ctx, "nie udalo sie pobrac konta!")
			return
		}

		if len(accounts) != 1 {

			if err != nil {
				views.ErrorPage(ctx, "nie udalo sie pobrac konta!")
				return
			}

		}

		records, err = db.Query("SELECT t.id, t.source_account_id, t.destination_account_id, t.amount, sa.holder_name, da.holder_name FROM account_transaction t LEFT JOIN account sa ON t.source_account_id = sa.id LEFT JOIN account da ON t.destination_account_id = da.id WHERE t.source_account_id = ? OR t.destination_account_id = ?", accountId, accountId)

		if err != nil {
			views.ErrorPage(ctx, "Nie udało się pobrać listy transakcji!"+err.Error())
			return
		}

		tx := make([]views.TransactionViewModel, 0)
		for records.Next() {
			var tvm views.TransactionViewModel

			err = records.Scan(&tvm.Id, &tvm.SourceAccountId, &tvm.DestinationAccountId, &tvm.Amount, &tvm.SourceAccountName, &tvm.DestinationAccountName)
			if err != nil {
				views.ErrorPage(ctx, "Nie udało się pobrać listy transakcji!"+err.Error())
				return
			}
			tx = append(tx, tvm)
		}

		balance, err := lib.GetAccountBalance(db, accountIdNum)
		ctx.HTML(http.StatusOK, "", views.AccountDetails(views.AccountDetailsViewModel{Account: accounts[0], Transactions: tx, Balance: balance}))
	}
}

func PostDashboardAccountDetailsDepositHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := lib.GetUser(ctx, db)
		if err != nil {
			views.ErrorPage(ctx, "Must be logged in!")
			return
		}

		accountIdStr := ctx.Param("accountId")
		amountStr := ctx.PostForm("amount")

		if accountIdStr == "" || amountStr == "" {
			views.ErrorPage(ctx, "account ID param and body field amount are required")
			return
		}

		accountId, err := strconv.Atoi(accountIdStr)
		if err != nil {
			views.ErrorPage(ctx, err.Error())
			return
		}
		amount, err := strconv.ParseFloat(amountStr, 32)
		if err != nil {
			views.ErrorPage(ctx, err.Error())
			return
		}

		tx, err := lib.CreateTransaction(db, -1, int64(accountId), amount)
		if err != nil {
			views.ErrorPage(ctx, err.Error())
			return
		}

		fmt.Sprintf("%+v\n", tx)

		ctx.Redirect(http.StatusFound, fmt.Sprintf("/dashboard/accounts/%d", accountId))
	}
}

func PostDashboardAccountDetailsWithdrawHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := lib.GetUser(ctx, db)
		if err != nil {
			views.ErrorPage(ctx, "Must be logged in!")
			return
		}

		accountIdStr := ctx.Param("accountId")
		amountStr := ctx.PostForm("amount")

		if accountIdStr == "" || amountStr == "" {
			views.ErrorPage(ctx, "account ID param and body field amount are required")
			return
		}

		accountId, err := strconv.Atoi(accountIdStr)
		if err != nil {
			views.ErrorPage(ctx, err.Error())
			return
		}
		amount, err := strconv.ParseFloat(amountStr, 32)
		if err != nil {
			views.ErrorPage(ctx, err.Error())
			return
		}

		tx, err := lib.CreateTransaction(db, int64(accountId), -1, amount)
		if err != nil {
			views.ErrorPage(ctx, err.Error())
			return
		}

		fmt.Sprintf("%+v\n", tx)

		ctx.Redirect(http.StatusFound, fmt.Sprintf("/dashboard/accounts/%d", accountId))
	}
}

func PostDashboardAccountDetailsCreateTransactionHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := lib.GetUser(ctx, db)
		if err != nil {
			views.ErrorPage(ctx, "Must be logged in!")
			return
		}

		sourceAccountIdStr := ctx.Param("accountId")
		destinationAccountIdStr := ctx.PostForm("destination-account-id")
		amountStr := ctx.PostForm("amount")

		if sourceAccountIdStr == "" || destinationAccountIdStr == "" || amountStr == "" {
			views.ErrorPage(ctx, "account ID param and body fields destination-account-id and amount are required")
			return
		}

		sourceAccountId, err := strconv.Atoi(sourceAccountIdStr)
		if err != nil {
			views.ErrorPage(ctx, err.Error())
			return
		}
		destinationAccountId, err := strconv.Atoi(destinationAccountIdStr)
		if err != nil {
			views.ErrorPage(ctx, err.Error())
			return
		}
		amount, err := strconv.ParseFloat(amountStr, 32)
		if err != nil {
			views.ErrorPage(ctx, err.Error())
			return
		}

		tx, err := lib.CreateTransaction(db, int64(sourceAccountId), int64(destinationAccountId), amount)
		if err != nil {
			views.ErrorPage(ctx, err.Error())
			return
		}

		fmt.Sprintf("%+v\n", tx)

		ctx.Redirect(http.StatusFound, fmt.Sprintf("/dashboard/accounts/%d", sourceAccountId))
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

func TransactionsHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := lib.GetUser(ctx, db)
		if err != nil {
			views.ErrorPage(ctx, "Must be logged in!")
			return
		}

		records, err := db.Query("SELECT t.id, t.source_account_id, t.destination_account_id, t.amount, sa.holder_name, da.holder_name FROM account_transaction t LEFT JOIN account sa ON t.source_account_id = sa.id LEFT JOIN account da ON t.destination_account_id = da.id")

		if err != nil {
			views.ErrorPage(ctx, "Nie udało się pobrać listy transakcji!"+err.Error())
			return
		}

		tx := make([]views.TransactionViewModel, 0)
		for records.Next() {
			var tvm views.TransactionViewModel

			err = records.Scan(&tvm.Id, &tvm.SourceAccountId, &tvm.DestinationAccountId, &tvm.Amount, &tvm.SourceAccountName, &tvm.DestinationAccountName)
			if err != nil {
				views.ErrorPage(ctx, "Nie udało się pobrać listy transakcji!"+err.Error())
				return
			}
			tx = append(tx, tvm)
		}

		ctx.HTML(http.StatusOK, "", views.Transactions(tx))
	}

}
