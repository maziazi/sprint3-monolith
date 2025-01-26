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
	ErrPhoneNotFound      = errors.New("phone not found")
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
		Email:    &email,
		Password: string(hashedPassword),
		Id:       userID,
	}, nil
}

func RegisterUserPhone(phone, password string) (*model.User, error) {
	db := database.GetDBPool()

	// Check if phone exists
	var existingUser model.User
	err := db.QueryRow(context.Background(), "SELECT phone FROM public.user WHERE phone = $1", phone).Scan(&existingUser.Phone)
	if err == nil {
		return nil, ErrPhoneAlreadyExists
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
		`INSERT INTO public.user (phone, password, "createdAt") 
         VALUES ($1, $2, $3) 
         RETURNING "userId"`,
		phone, string(hashedPassword), time.Now(),
	).Scan(&userID)

	if err != nil {
		return nil, fmt.Errorf("failed to register user: %v", err)
	}

	_, err = db.Exec(context.Background(), `INSERT INTO "userProfile" (phone, "userId") VALUES ($1, $2)`,
		phone, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user profile: %v", err)
	}

	return &model.User{
		Phone:    &phone,
		Password: string(hashedPassword),
		Id:       userID,
	}, nil
}

func AuthenticateEmail(email, password string) (*model.User, error) {
	db := database.GetDBPool()
	var user model.User

	// Retrieve user by email, gunakan pointer untuk phone agar bisa menangani NULL
	err := db.QueryRow(context.Background(), `SELECT "userId", email, password, phone FROM public.user WHERE email = $1`, email).
		Scan(&user.Id, &user.Email, &user.Password, &user.Phone)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrEmailNotFound
	} else if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidPassword
	}

	return &user, nil
}

func AuthenticatePhone(phone, password string) (*model.User, error) {
	db := database.GetDBPool()
	var user model.User

	// Retrieve user by email, gunakan pointer untuk phone agar bisa menangani NULL
	err := db.QueryRow(context.Background(), `SELECT "userId", email, password, phone FROM public.user WHERE phone = $1`, phone).
		Scan(&user.Id, &user.Email, &user.Password, &user.Phone)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrEmailNotFound
	} else if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidPassword
	}

	return &user, nil
}
