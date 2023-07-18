package third

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func LoginGithub(c *gin.Context) {
	path := fmt.Sprintf("https://github.com/login/oauth/authorize?scope=user&client_id=%s&state=%s&redirect_uri=%s", githubConfig.client_id, authState, githubConfig.redirect_uri)

	c.Redirect(http.StatusMovedPermanently, path)
}

func CallbackGithub(c *gin.Context) {

	code := c.Query("code")

	client := &http.Client{}

	data := url.Values{}
	data.Set("client_id", githubConfig.client_id)
	data.Set("client_secret", githubConfig.client_secret)
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

	var userInfo UserInfo

	infoErr := json.Unmarshal(bodyBytesUserInfo, &userInfo)

	if infoErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "04",
			"message": "解析失败",
		})
		return
	}

	c.SetCookie("username", userInfo.Login, MAX_AGE, "/", Cookie_Host, false, true)
	c.SetCookie("avatar_url", userInfo.Avatar_url, MAX_AGE, "/", Cookie_Host, false, true)

	c.Redirect(http.StatusMovedPermanently, homePage)
}
