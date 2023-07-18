package main

import (
	"net/http"

	"third-auth-go/third"

	"third-auth-go/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Cookie"},
		AllowCredentials: true,
		AllowOrigins:     []string{"http://localhost:3001"},
	}))

	r.GET("/login/github", third.LoginGithub)
	r.GET("/callback/github", third.CallbackGithub)
	r.GET("/login/gitee", third.LoginGitee)
	r.GET("/callback/gitee", third.CallbackGitee)
	r.GET("/user", User)
	r.GET("/logout", Logout)

	r.Run(":5501")
}

func User(c *gin.Context) {
	username, _ := c.Cookie("username")
	avatar_url, _ := c.Cookie("avatar_url")

	if username == "" || avatar_url == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    "0",
			"message": "success",
			"data": map[string]interface{}{
				"access": false,
			},
		})
		return
	}

	data := map[string]interface{}{
		"username":   username,
		"avatar_url": avatar_url,
		"access":     utils.Includes(third.WhiteList, username),
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "0",
		"message": "success",
		"data":    data,
	})
}

func Logout(c *gin.Context) {
	c.SetCookie("username", "", -1, "/", "", false, true)
	c.SetCookie("avatar_url", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"code":    "0",
		"message": "success",
		"data":    true,
	})
}
