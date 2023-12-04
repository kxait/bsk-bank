package handlers

import (
	"bsk-bank/lib"
	"bsk-bank/views"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
