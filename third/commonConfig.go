package third

var homePage = "http://localhost:3001"
var authState = "big_deal"
var MAX_AGE = 60 * 60 * 24
var Cookie_Host = "localhost:3001."

type Config struct {
	client_id     string
	client_secret string
	redirect_uri  string
}

var githubConfig = Config{
	client_id:     "95a6cf45f087a9b6dbe2",
	client_secret: "7efaeec759e2eede82a05d37dd1c33196443c9f8",
	redirect_uri:  "http://localhost:3001/callback/github",
}

var giteeConfig = Config{
	client_id:     "c8632d627cd786d4a76386aaa3ac4abe656a178991b9b6fd21334e0ce0c68ae7",
	client_secret: "3a96bea00210a952face9bfc7b7fe1dfaba815bf5e6e0041523242c053b8a6ed",
	redirect_uri:  "http://localhost:3001/callback/gitee",
}

type UserInfo struct {
	Login      string `json:"login"`
	Avatar_url string `json:"avatar_url"`
}
