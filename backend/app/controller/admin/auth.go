package admin

import (
	"os"
	"amidesk-api-server/app/form/admin"
	"amidesk-api-server/app/model"
	"amidesk-api-server/config"
	"amidesk-api-server/helper/captcha"
	"amidesk-api-server/helper/security"
	"amidesk-api-server/util"
	"strconv"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/pquerna/otp/totp"
)

type AuthController struct {
	basicController
	Cfg *config.ServerConfig
}

func (c *AuthController) PostAuthLogin() mvc.Result {
	var loginForm admin.LoginForm
	err := c.Ctx.ReadJSON(&loginForm)
	if err != nil {
		return c.Error(nil, err.Error())
	}
	if !security.Allow(c.Cfg.Security, "admin", loginForm.Username) {
		return c.Error(nil, "TooManyLoginAttempts")
	}

	if os.Getenv("E2E_SKIP_CAPTCHA") != "true" && !captcha.VerifyCode(loginForm.CaptchaId, loginForm.Code) {
		security.RecordFailure(c.Cfg.Security, "admin", loginForm.Username)
		return c.Error(nil, "CaptchaError")
	}

	var user model.User
	get, err := c.Db.Where("username = ? and is_admin = 1", loginForm.Username).Get(&user)
	if err != nil {
		return c.Error(nil, err.Error())
	}

	if !get {
		security.RecordFailure(c.Cfg.Security, "admin", loginForm.Username)
		return c.Error(nil, "AuthenticationFailed")
	}

	if !util.PasswordVerify(loginForm.Password, user.Password) {
		security.RecordFailure(c.Cfg.Security, "admin", loginForm.Username)
		return c.Error(nil, "AuthenticationFailed")
	}

	if c.Cfg.Security.RequireAdminTOTP && (user.TwoFactorAuthSecret == "" || !totp.Validate(loginForm.TfaCode, user.TwoFactorAuthSecret)) {
		security.RecordFailure(c.Cfg.Security, "admin", loginForm.Username)
		return c.Error(nil, "AuthenticationFailed")
	}
	security.Clear("admin", loginForm.Username)

	// make other tokens expired
	_, _ = c.Db.Where("user_id = ? and status = 1 and is_admin = 1", user.Id).Cols("status").Update(&model.AuthToken{
		Status: 0,
	})

	signStr := strconv.Itoa(user.Id) + user.Username + time.Now().String()
	token := util.HmacSha256(signStr, c.Cfg.SignKey)
	expired := 2 * time.Hour // 2 hours

	authToken := &model.AuthToken{
		UserId:  user.Id,
		Token:   token,
		Expired: time.Now().Add(expired),
		IsAdmin: true,
		Status:  1,
	}

	_, err = c.Db.Insert(authToken)
	if err != nil {
		return c.Error(nil, err.Error())
	}

	return c.Success(iris.Map{
		"token": token,
	}, "ok")
}

func (c *AuthController) GetAuthCaptcha() mvc.Result {
	id, img := captcha.CreateCaptcha()
	return c.Success(iris.Map{
		"id":  id,
		"img": img,
	}, "ok")
}
