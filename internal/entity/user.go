package entity

type User struct {
	BaseEntity
	Username string `gorm:"uniqueIndex" json:"username"`
	Password string `json:"-"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}