package app

import (
	"testing"

	"github.com/appleboy/gofight"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	c := gofight.New()
	var r *gin.Engine
	r = Setup("test")

	c.POST("/register").
		SetForm(gofight.H{
			"Username":    "username",
			"Password":    "12345678",
			"DisplayName": "woo",
		}).
		Run(r, func(res gofight.HTTPResponse, req gofight.HTTPRequest) {
			user := getUser("username")
			assert.NotNil(t, user)
			assert.Equal(t, "username", user.Username)
			assert.Equal(t, "12345678", user.Password)
			assert.Equal(t, "woo", user.DisplayName)
		})
}
