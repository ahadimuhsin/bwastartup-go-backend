package handler

import (
	"bwastartup/helper"
	"bwastartup/transaction"
	"bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *transactionHandler {
	return &transactionHandler{service}
}

func (h *transactionHandler) GetCampaignTransactions(c *gin.Context) {
	//parameter di url
	//tangkap parameter, mapping ke input struct
	var input transaction.GetCampaignTransactionsInput

	err := c.ShouldBindUri(&input)

	if err != nil {
		response := helper.APIResponse("Input URI Error", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//ambil current user dari jwt/handler
	//ambil data user dari Context gin, dari auth Middleware
	currentUser := c.MustGet("currentUser").(user.User)
	//set nilai idnya ke variabel
	input.User = currentUser

	//panggil service, input struct sebagai parameter
	//service : panggil repository
	//repo: mencari data transaction suatu campaign
	transactions, err := h.service.GetTransactionByCampaignID(input)

	if err != nil {
		response := helper.APIResponse("Failed to get list of campaign's transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign's Transactions Detail", http.StatusOK, "success", transaction.FormatCampaignTransactions(transactions))
	c.JSON(http.StatusOK, response)
}


func (h *transactionHandler) GetUserTransactions(c *gin.Context){
	// handler: ambil nilai user dari jwt
	//ambil current user dari jwt/handler
	//ambil data user dari Context gin, dari auth Middleware
	currentUser := c.MustGet("currentUser").(user.User)
	//set nilai idnya ke variabel
	userID := currentUser.ID
	//service
	//repo -> ambil data transaction (preload campaign)
	transactions, err := h.service.GetTransactionByUserID(userID)

	if err != nil{
		response := helper.APIResponse("Failed to get list of user's transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	

	response := helper.APIResponse("User's Transactions", http.StatusOK, "success", transaction.FormatUserTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) CreateTransaction(c *gin.Context){

	//input nominal dari user
	//tangkap input, mapping ke input struct
	var input transaction.CreateTransactionInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}
		//response menggunakna helper
		response := helper.APIResponse("Failed To Create Transaction", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	//ambil current user dari jwt/handler
	//ambil data user dari Context gin, dari auth Middleware
	currentUser := c.MustGet("currentUser").(user.User)
	//set nilai idnya ke variabel
	input.User = currentUser

	//panggil service buat transaksi, manggil sistem midtrans
	//panggi repo, create new transaction
	newTransaction, err := h.service.CreateTransaction(input)

	if err != nil {
		response := helper.APIResponse("Failed to save transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Transaction Created", http.StatusOK, "success", transaction.FormatTransaction(newTransaction))
	c.JSON(http.StatusOK, response)
}


