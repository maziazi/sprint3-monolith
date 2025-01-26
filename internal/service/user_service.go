package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
	"sprint3/internal/model"
	"sprint3/pkg/database"
	"time"
)

var (
	ErrEmailNotFound      = errors.New("email not found")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrPhoneAlreadyExists = errors.New("phone already exists")
)

func RegisterUserEmail(email, password string) (*model.User, error) {
	db := database.GetDBPool()

	// Check if email exists
	var existingUser model.User
	err := db.QueryRow(context.Background(), "SELECT email FROM public.user WHERE email = $1", email).Scan(&existingUser.Email)
	if err == nil {
		return nil, ErrEmailAlreadyExists
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("database error: %v", err)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// Insert user into database and get the generated ID
	var userID uint
	err = db.QueryRow(context.Background(),
		`INSERT INTO public.user (email, password, "createdAt") 
         VALUES ($1, $2, $3) 
         RETURNING "userId"`,
		email, string(hashedPassword), time.Now(),
	).Scan(&userID)

	if err != nil {
		return nil, fmt.Errorf("failed to register user: %v", err)
	}

	_, err = db.Exec(context.Background(), `INSERT INTO "userProfile" (email, "userId") VALUES ($1, $2)`,
		email, userID)
	return &model.User{
		Email:    email,
		Password: string(hashedPassword),
		Id:       userID,
	}, nil
}
