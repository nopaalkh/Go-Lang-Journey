package middlewares

import (
	"net/http"

	"belajar-auth/config"
	"belajar-auth/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil Cookie
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// 2. Cek apakah tokennya valid
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte("rahasia-kita"), nil
		})

		if err != nil || !token.Valid {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// 3. Ambil Data User dari Database
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Ambil ID dari token (format float64)
			if floatId, ok := claims["sub"].(float64); ok {

				var user models.User
				if err := config.DB.First(&user, int64(floatId)).Error; err != nil {
					c.SetCookie("token", "", -1, "/", "localhost", false, true)
					c.Redirect(http.StatusFound, "/login")
					c.Abort()
					return
				}

				c.Set("user", user)

			} else {
				c.Redirect(http.StatusFound, "/login")
				c.Abort()
				return
			}
		} else {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		c.Next()
	}
}
