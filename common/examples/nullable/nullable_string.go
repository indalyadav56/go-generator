package main

import (
	"common/types"
	"fmt"
)

// Example with struct
type User struct {
	ID         int64                `json:"id" db:"id"`
	FirstName  types.NullableString `json:"first_name" db:"first_name"`
	MiddleName types.NullableString `json:"middle_name" db:"last_name"`
	LastName   types.NullableString `json:"last_name" db:"last_name"`
	NickName   types.NullableString `json:"nick_name" db:"last_name"`
}

func main() {

	user1 := new(User)
	// Method 1: Using constructors
	firstName := "John"
	user1.FirstName = types.NewNullableString(&firstName)

	// Method 2: Using FromString
	user1.LastName = types.FromString("Doe")

	// Method 3: Using SetValid
	var user User
	user.MiddleName.SetValid("Middle")

	// Method 4: Setting null/invalid
	user.NickName.SetInvalid()

	// Method 5: Getting value safely
	name := user.FirstName.GetValue()

	fmt.Println("name", name)
}

// Example with request DTO
type CreateUserRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
}

func (r *CreateUserRequest) ToUser() User {
	return User{
		FirstName: types.NewNullableString(r.FirstName),
		LastName:  types.NewNullableString(r.LastName),
	}
}
