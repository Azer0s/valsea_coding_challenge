package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"valsea_coding_challenge/domain/transactional"
	"valsea_coding_challenge/service"
)

func AccountController(r *gin.Engine, log *zap.Logger, userService service.UserService) {
	r.POST("/accounts", func(c *gin.Context) {
		log.Debug("Received request to create account")

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

		user, err := userService.CreateUser(req)
		if err != nil {
			log.Error("Failed to create account", zap.Error(err))
			c.JSON(500, gin.H{"error": "Failed to create account"})
			return
		}

		c.JSON(201, user)
	})

	r.GET("/accounts/:id", func(c *gin.Context) {
		log.Info("Received request to get account by id")

		id := c.Param("id")
		user, err := userService.GetUserById(id)
		if err != nil {
			log.Error("Failed to get account by id", zap.Error(err))
			c.JSON(404, gin.H{"error": "Failed to get account by id"})
			return
		}

		c.JSON(200, user)
	})

	r.GET("/accounts", func(c *gin.Context) {
		log.Info("Received request to get all accounts")

		users, err := userService.GetAllUsers()
		if err != nil {
			log.Error("Failed to get all accounts", zap.Error(err))
			c.JSON(500, gin.H{"error": "Failed to get all accounts"})
			return
		}

		c.JSON(200, users)
	})
}
