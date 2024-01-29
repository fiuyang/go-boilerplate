package middleware

import (
	"fmt"
	"gin-boilerplate/config"
	"gin-boilerplate/helper"
	"gin-boilerplate/repository"
	"gin-boilerplate/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func JwtMiddleware(userRepository repository.UsersRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var token string
		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)
		
		if len(fields) != 0 && fields[0] == "Bearer" {
			token = fields[1]
		}
		
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Unauthorized"})
			return
		}
		
		if utils.IsTokenBlacklisted(token) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Token has been blacklisted"})
			return
		}

		config, _ := config.LoadConfig(".")
		sub, err := utils.ValidateToken(token, config.TokenSecret)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		
		id, err_id := strconv.Atoi(fmt.Sprint(sub))
		helper.ErrorPanic(err_id)
		result, err := userRepository.FindById(id)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
			return
		}
		
		ctx.Set("currentUser", result.Username)
		ctx.Next()
		
	}
}
