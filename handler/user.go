package handler

import (
	dblayer "filestore-server/db"
	"filestore-server/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	// 用于加密的盐值(自定义)
	pwdSalt = "*#890"
)

// SignupHandler : 处理用户注册请求
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(data)
		return
	}
	_ = r.ParseForm()

	username := r.Form.Get("username")
	passwd := r.Form.Get("password")

	if len(username) < 3 || len(passwd) < 5 {
		_, _ = w.Write([]byte("Invalid parameter"))
		return
	}

	// 对密码进行加盐及取Sha1值加密
	encPasswd := util.Sha1([]byte(passwd + pwdSalt))
	// 将用户信息注册到用户表中
	suc := dblayer.UserSignup(username, encPasswd)
	if suc {
		_, _ = w.Write([]byte("SUCCESS"))
	} else {
		_, _ = w.Write([]byte("FAILED"))
	}
}

// SignInHandler : 登录接口
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	encPasswd := util.Sha1([]byte(password + pwdSalt))
	// 1. 校验用户名及密码
	pwdChecked := dblayer.UserSignin(username, encPasswd)
	if !pwdChecked {
		_, _ = w.Write([]byte("FAILED"))
		return
	}

	// 2. 生成访问凭证(token)
	token := GenToken(username)
	upRes := dblayer.UpdateToken(username, token)
	if !upRes {
		_, _ = w.Write([]byte("FAILED"))
		return
	}

	// 3. 登录成功后重定向到首页
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	_, _ = w.Write(resp.JSONBytes())
}

// UserInfoHandler ： 查询用户信息
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求参数
	_ = r.ParseForm()
	username := r.Form.Get("username")
	token := r.Form.Get("token")

	// 2. 验证token是否有效
	isValidToken := IsTokenValid(token)
	if !isValidToken {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// 3. 查询用户信息
	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// 4. 组装并且响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	_, _ = w.Write(resp.JSONBytes())
}

// GenToken : 生成token
func GenToken(username string) string {
	// 40位字符:md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

// IsTokenValid : token是否有效
func IsTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}
	return true
}
