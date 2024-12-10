package main

import (
	"fmt"
	"main/backend"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func generateToken(c *gin.Context) {
	var authRequest backend.AuthRequest
	if err := c.ShouldBindJSON(&authRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid json",
		})
		return
	}

	var user backend.User
	backend.DB.First(&user, "username = ?", authRequest.Username)

	if user.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})
		return
	}

	userItem := &backend.UserItem{
		Id:       user.Id,
		Username: user.Username,
		Password: user.Password,
		RoleId:   user.Role.Id,
		Role: backend.RoleItem{
			Id:   user.Role.Id,
			Name: user.Role.Name,
		},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	token, refreshToken, err := backend.GenerateAuthRefreshTokens(userItem)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create accessToken",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  token,
		"refreshToken": refreshToken,
	})
}

func refreshToken(c *gin.Context) {
	var refreshTokenRequest backend.RefreshTokenRequest
	if err := c.ShouldBindJSON(&refreshTokenRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid json request",
		})
		return
	}

	claims := &backend.UserItem{}
	token, err := jwt.ParseWithClaims(refreshTokenRequest.RefreshToken, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(backend.SECRETKEY), nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Can't parse jwt claims: %s", err.Error()),
		})
		return
	}

	if !token.Valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid token",
		})
		return
	}

	newToken, err := backend.CreateToken(claims, backend.ACCESSTOKENTIME, &claims.Guid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error to refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": newToken,
	})
}

func main() {
	err := backend.DatabaseConnect()

	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.POST("/auth/token", generateToken)
	r.POST("/auth/refresh", refreshToken)
	r.Run()
}
