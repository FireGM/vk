package client

type Error struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_msg"`
}

func (e Error) Error() string {
	return e.Message
}

type ErrorResponse struct {
	Err  Error   `json:"error"`
	Errs []Error `json:"execute_errors"`
}

type ACError struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_msg"`
}

type ResponseCaptcha struct {
	Error struct {
		Captcha_sid string
		Captcha_img string
	}
}

type ResponseAuthLoginPass struct {
	AccessToken string `json:"access_token"`
	UserId      int    `json:"user_id"`
}

type userCheckResponse struct {
	Response []struct {
		Id int
	}
}