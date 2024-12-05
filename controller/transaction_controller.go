package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"valsea_coding_challenge/domain/dto"
	"valsea_coding_challenge/domain/enum"
	"valsea_coding_challenge/domain/transactional"
	"valsea_coding_challenge/service"
)

func TransactionController(r *gin.Engine, userService service.UserService, transactionService service.TransactionService, log *zap.Logger) {
	r.GET("/accounts/:id/transactions", func(c *gin.Context) {
		log.Info("Received request to get transactions")

		userId := c.Param("id")
		_, err := userService.GetUserById(userId)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to get transactions"})
			return
		}

		transactions, err := transactionService.GetTransactions(userId)
		if err != nil {
			return
		}

		c.JSON(200, transactions)
	})

	r.POST("/accounts/:id/transactions", func(c *gin.Context) {
		log.Info("Received request to create transaction")

		userId := c.Param("id")
		_, err := userService.GetUserById(userId)
		if err != nil {
			log.Error("Failed to get user", zap.Error(err))
			c.JSON(404, gin.H{"error": "Failed to get user"})
			return
		}

		var request transactional.CreateTransactionRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			log.Error("Failed to bind request", zap.Error(err))
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		var transaction *dto.TransactionDTO
		switch request.Type {
		case enum.TransactionTypeDeposit:
			transaction, err = transactionService.Deposit(userId, request.Amount)
		case enum.TransactionTypeWithdrawal:
			transaction, err = transactionService.Withdraw(userId, request.Amount)
		}

		if err != nil {
			log.Error("Failed to create transaction", zap.Error(err))
			c.JSON(500, gin.H{"error": "Failed to create transaction"})
			return
		}

		c.JSON(200, transaction)
	})

	r.POST("/transfer", func(c *gin.Context) {
		log.Info("Received request to transfer")

		var request transactional.CreateTransferRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			log.Error("Failed to bind request", zap.Error(err))
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		_, err := userService.GetUserById(request.FromUserId)
		if err != nil {
			log.Error("Failed to get user", zap.Error(err))
			c.JSON(404, gin.H{"error": "Failed to get user"})
			return
		}

		_, err = userService.GetUserById(request.ToUserId)
		if err != nil {
			log.Error("Failed to get user", zap.Error(err))
			c.JSON(404, gin.H{"error": "Failed to get user"})
			return
		}

		transaction, err := transactionService.Transfer(request.FromUserId, request.ToUserId, request.Amount)
		if err != nil {
			log.Error("Failed to create transaction", zap.Error(err))
			c.JSON(500, gin.H{"error": "Failed to create transaction"})
			return
		}

		c.JSON(200, transaction)
	})
}
