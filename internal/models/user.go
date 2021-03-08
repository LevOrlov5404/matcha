package models

type (
	UserToCreate struct {
		Email     string `json:"email" binding:"required"`
		Username  string `json:"username" binding:"required"`
		FirstName string `json:"firstName" binding:"required"`
		LastName  string `json:"lastName" binding:"required"`
		Password  string `json:"password" binding:"required"`
	}

	UserToSignIn struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	User struct {
		ID        uint64 `json:"id" db:"id"`
		Email     string `json:"email" db:"email"`
		Username  string `json:"username" db:"username"`
		FirstName string `json:"firstName" db:"first_name"`
		LastName  string `json:"lastName" db:"last_name"`
		Password  string `json:"-" db:"password"`
	}
)
