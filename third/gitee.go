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

func LoginGitee(c *gin.Context) {
	path := fmt.Sprintf("https://gitee.com/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code", giteeConfig.client_id, giteeConfig.redirect_uri)

	c.Redirect(http.StatusMovedPermanently, path)
}

func CallbackGitee(c *gin.Context) {
	code := c.Query("code")

	client := &http.Client{}

	data := url.Values{}

	data.Set("grant_type", "authorization_code")
	data.Set("client_id", giteeConfig.client_id)
	data.Set("client_secret", giteeConfig.client_secret)
	data.Set("redirect_uri", giteeConfig.redirect_uri)
	data.Set("code", code)

	req, err := http.NewRequest(http.MethodPost, "https://gitee.com/oauth/token", strings.NewReader(data.Encode()))

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

	fmt.Println("result", result)

	access_token := result["access_token"]

	infoPath := fmt.Sprintf("https://gitee.com/api/v5/user?access_token=%s", access_token)

	userInfoReq, _ := http.NewRequest(http.MethodGet, infoPath, nil)

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
