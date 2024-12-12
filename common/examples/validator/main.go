package main

// import (
// 	"fmt"
// 	"paydoh-common/pkg/validator"
// 	"time"
// )

// // User represents a user entity with validation tags
// type User struct {
// 	ID        string    `json:"id" validate:"required,uuid"`
// 	Email     string    `json:"email" validate:"required,email"`
// 	Password  string    `json:"password" validate:"required,password"`
// 	Phone     string    `json:"phone" validate:"required,phone"`
// 	CreatedAt time.Time `json:"created_at" validate:"required"`
// }

// func main() {
// 	// Create a new validator instance
// 	v, err := validator.NewCustomValidator()
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Example usage with struct
// 	user := User{
// 		ID:        "invalid-uuid",
// 		Email:     "invalid-email",
// 		Password:  "weak",
// 		Phone:     "invalid",
// 		CreatedAt: time.Now(),
// 	}

// 	if err := v.Validate(user); err != nil {
// 		if validationErrors, ok := err.(validator.ValidationErrors); ok {
// 			for _, e := range validationErrors {
// 				fmt.Printf("Field: %s, Error: %s\n", e.Field, e.Message)
// 			}
// 		}
// 	}

// 	// Example of single field validation
// 	email := "invalid-email"
// 	if err := v.ValidateField(email, "email"); err != nil {
// 		fmt.Printf("Email validation error: %v\n", err)
// 	}

// 	// Example of JSON validation
// 	jsonData := []byte(`{
// 		"id": "123e4567-e89b-12d3-a456-426614174000",
// 		"email": "user@example.com",
// 		"password": "StrongPass1!",
// 		"phone": "+1234567890",
// 		"created_at": "2024-10-26T12:00:00Z"
// 	}`)

// 	var newUser User
// 	if err := v.ValidateJSON(jsonData, &newUser); err != nil {
// 		fmt.Printf("JSON validation error: %v\n", err)
// 	}
// }
