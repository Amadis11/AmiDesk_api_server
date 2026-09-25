package admin

type LoginForm struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	TfaCode   string `json:"tfaCode"`
	Code      string `json:"code"`
	CaptchaId string `json:"captchaId"`
}
