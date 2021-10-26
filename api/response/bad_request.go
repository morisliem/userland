package response

type Bad_Request struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Bad_request(s string) Bad_Request {
	res := Bad_Request{
		Code:    0,
		Message: s,
	}

	return res
}
