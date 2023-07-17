package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var client_id = "95a6cf45f087a9b6dbe2"
var client_secret = "7efaeec759e2eede82a05d37dd1c33196443c9f8"
var redirect_uri = "http://localhost:3001/callback/github"
var homePage = "http://localhost:3001"
var authState = "big_deal"

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Cookie"},
		AllowCredentials: true,
		AllowOrigins:     []string{"http://localhost:3001"},
	}))

	r.GET("/login/github", loginGithub)
	r.GET("/callback/github", callbackGithub)
	r.GET("/user", User)

	r.Run(":5501")
}

func loginGithub(c *gin.Context) {
	path := fmt.Sprintf("https://github.com/login/oauth/authorize?scope=user&client_id=%s&state=%s&redirect_uri=%s", client_id, authState, redirect_uri)

	c.Redirect(http.StatusMovedPermanently, path)
}

func callbackGithub(c *gin.Context) {

	code := c.Query("code")

	client := &http.Client{}

	data := url.Values{}
	data.Set("client_id", client_id)
	data.Set("client_secret", client_secret)
	data.Set("code", code)

	req, err := http.NewRequest(http.MethodPost, "https://github.com/login/oauth/access_token", strings.NewReader(data.Encode()))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "01",
			"message": "构建请求失败",
		})
		return
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "02",
			"message": "发送请求失败",
		})
		return
	}

	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var result map[string]interface{}

	json.Unmarshal(bodyBytes, &result)

	if err, ok := result["error"]; ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "02",
			"message": err,
		})
		return
	}

	access_token := result["access_token"]
	token_type := result["token_type"]
	authToken := fmt.Sprintf("%s %s", token_type, access_token)

	userInfoReq, _ := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	userInfoReq.Header.Add("Authorization", authToken)
	respUserInfo, errUserInfo := client.Do(userInfoReq)

	if errUserInfo != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "03",
			"message": "获取用户信息失败：" + errUserInfo.Error(),
		})
		return
	}
	defer respUserInfo.Body.Close()
	bodyBytesUserInfo, _ := ioutil.ReadAll(respUserInfo.Body)
	var userInfo map[string]interface{}
	json.Unmarshal(bodyBytesUserInfo, &userInfo)

	fmt.Println("userInfo: ", userInfo)

	c.SetCookie("username", "dilireba", 60*60*24, "/", "localhost:3001.", false, true)
	c.SetCookie("age", "18", 60*60*24, "/", "localhost:3001.", false, true)

	c.Redirect(http.StatusMovedPermanently, homePage)
}

func User(c *gin.Context) {
	username, _ := c.Cookie("username")
	age, _ := c.Cookie("age")

	fmt.Println("username: ", username)
	fmt.Println("age: ", age)

	c.JSON(http.StatusOK, gin.H{
		"code":     "0",
		"message":  "success",
		"username": username,
		"age":      age,
	})
}
