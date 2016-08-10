package anticaptcha

type ResponseCaptcha struct {
	Error struct {
		Captcha_sid string
		Captcha_img string
	}
}