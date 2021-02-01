package controllers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/kazetora/evermos-assignment/problem_1_ecommerce/helpers"
	"github.com/kazetora/evermos-assignment/problem_1_ecommerce/models"
	"github.com/kazetora/evermos-assignment/problem_1_ecommerce/storage"
)

type TransactionResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func GetTransactionStatus(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	transaction, err := storage.GetTransaction(id)
	if err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := TransactionResponse{
		Status: transaction.Status,
	}

	switch transaction.Status {
	case models.TransactionStatusUnprocessed:
		w.WriteHeader(http.StatusNoContent)
		return
	case models.TransactionStatusSuccess:
		resp.Message = "Transaction has succeeded"
	case models.TransactionStatusError:
		resp.Message = "Transaction has failed"
	}

	helpers.RespondJSON(w, http.StatusOK, resp)
}
