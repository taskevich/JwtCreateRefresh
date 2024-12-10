package backend

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type RoleItem struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type UserItem struct {
	Id        uint      `json:"id"`
	Guid      string    `json:"guid"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	RoleId    uint      `json:"roleId"`
	Role      RoleItem  `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	jwt.StandardClaims
}
