package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"valsea_coding_challenge/domain/transactional"
	"valsea_coding_challenge/service"
)

func AccountController(r *gin.Engine, log *zap.Logger, accountService service.AccountService) {
	r.POST("/accounts", func(c *gin.Context) {
		log.Info("Received request to create account")

		reqBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Error("Failed to read request body", zap.Error(err))
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		var req transactional.CreateUserRequest
		err = json.Unmarshal(reqBody, &req)
		if err != nil {
			log.Error("Failed to unmarshal request body", zap.Error(err))
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		user, err := accountService.CreateAccount(req)
		if err != nil {
			log.Error("Failed to create account", zap.Error(err))
			c.JSON(500, gin.H{"error": "Failed to create account"})
			return
		}

		c.JSON(201, user)
	})

	//TODO GET /accounts/:id
	//TODO GET /accounts
	//TODO POST /accounts/:id/transactions
	//TODO GET /accounts/:id/transactions
	//TODO POST /transfer
}
