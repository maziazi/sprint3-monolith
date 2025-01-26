package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sprint3/internal/middleware"
	"sprint3/internal/service"
)

type AuthRequestEmail struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

type AuthRequestPhone struct {
	Phone    string `json:"phone" binding:"required,phone"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

func RegisterUserEmail(c *gin.Context) {
	log.Println("Handler RegisterUser hit")
	var req AuthRequestEmail
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON binding failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Input validated: %+v", req)
	user, err := service.RegisterUserEmail(req.Email, req.Password)
	if err != nil {
		log.Printf("Service error: %v", err)
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	token, _ := middleware.GenerateToken(user.Email, user.Id)
	if err != nil {
		log.Printf("Token generation failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	log.Println("User registered successfully")
	c.JSON(http.StatusCreated, gin.H{"email": user.Email, "phone": user.Phone, "token": token})
}
