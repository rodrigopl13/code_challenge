package entities

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name" form:"first_name" binding:"required"`
	LastName  string `json:"last_name" form:"last_name" binding:"required"`
	UserName  string `json:"user_name" form:"user_name" binding:"required"`
	Password  string `json:"password,omitempty" form:"password" binding:"required"`
}
