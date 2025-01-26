package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"regexp"
	"sprint3/internal/middleware"
	"sprint3/internal/service"
)

type AuthRequestEmail struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

type AuthRequestPhone struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

func RegisterUserEmail(c *gin.Context) {
	log.Println("Handler RegisterUserEmail hit")
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

	token, err := middleware.GenerateToken(user.Email, user.Phone, user.Id)
	if err != nil {
		log.Printf("Token generation failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	log.Println("User registered successfully")
	// Prepare response
	response := gin.H{
		"email": user.Email,
		"token": token,
	}

	// If phone is nil, set it to an empty string
	if user.Phone == nil {
		response["phone"] = ""
	} else {
		response["phone"] = *user.Phone
	}

	// Return response
	c.JSON(http.StatusOK, response)
}
func LoginUserEmail(c *gin.Context) {
	log.Println("Handler LoginUserEmail hit")
	var req AuthRequestEmail
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON binding failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Input validated: %+v", req)
	user, err := service.AuthenticateEmail(req.Email, req.Password)
	if err != nil {
		log.Printf("Authentication error: %v", err)
		if errors.Is(err, service.ErrEmailNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Email not found"})
		} else if errors.Is(err, service.ErrInvalidPassword) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	token, err := middleware.GenerateToken(user.Email, user.Phone, user.Id)
	if err != nil {
		log.Printf("Token generation failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	log.Println("User logged in successfully")
	// Prepare response
	response := gin.H{
		"email": user.Email,
		"token": token,
	}

	// If phone is nil, set it to an empty string
	if user.Phone == nil {
		response["phone"] = ""
	} else {
		response["phone"] = *user.Phone
	}

	// Return response
	c.JSON(http.StatusOK, response)
}

func RegisterUserPhone(c *gin.Context) {
	log.Println("Handler RegisterUserPhone hit")
	var req AuthRequestPhone
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON binding failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi tambahan untuk field `phone`
	if !isValidPhone(req.Phone) {
		log.Println("Phone validation failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone format. It must start with '+' followed by digits."})
		return
	}

	log.Printf("Input validated: %+v", req)
	user, err := service.RegisterUserPhone(req.Phone, req.Password)
	if err != nil {
		log.Printf("Service error: %v", err)
		if errors.Is(err, service.ErrPhoneAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "Phone already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	token, err := middleware.GenerateToken(user.Email, user.Phone, user.Id)
	if err != nil {
		log.Printf("Token generation failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	log.Println("User registered successfully")
	// Prepare response
	response := gin.H{
		"phone": user.Phone,
		"token": token,
	}

	// If phone is nil, set it to an empty string
	if user.Email == nil {
		response["email"] = ""
	} else {
		response["email"] = *user.Email
	}

	// Return response
	c.JSON(http.StatusOK, response)
}
func LoginUserPhone(c *gin.Context) {
	log.Println("Handler LoginUserPhone hit")
	var req AuthRequestPhone

	// Validasi tambahan untuk field `phone`

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON binding failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Input validated: %+v", req)
	user, err := service.AuthenticatePhone(req.Phone, req.Password)
	if err != nil {
		log.Printf("Authentication error: %v", err)
		if errors.Is(err, service.ErrPhoneNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Phone not found"})
		} else if errors.Is(err, service.ErrInvalidPassword) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}
	if !isValidPhone(req.Phone) {
		log.Println("Phone validation failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone format. It must start with '+' followed by digits."})
		return
	}
	token, err := middleware.GenerateToken(user.Email, user.Phone, user.Id)
	if err != nil {
		log.Printf("Token generation failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	log.Println("User logged in successfully")
	// Prepare response
	response := gin.H{
		"phone": user.Phone,
		"token": token,
	}

	// If phone is nil, set it to an empty string
	if user.Email == nil {
		response["email"] = ""
	} else {
		response["email"] = *user.Email
	}

	// Return response
	c.JSON(http.StatusOK, response)
}
func isValidPhone(phone string) bool {
	// Validasi nomor telepon: harus dimulai dengan "+" diikuti angka
	match, _ := regexp.MatchString(`^\+\d+$`, phone)
	return match
}
