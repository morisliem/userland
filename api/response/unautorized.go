package response

type Unautorized_Request struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Unautorized_request(s string) Unautorized_Request {
	res := Unautorized_Request{
		Code:    401,
		Message: s,
	}

	return res
}
