package response

type SuccessRes struct {
	Status bool `json:"status"`
}

func Success() SuccessRes {
	res := SuccessRes{
		Status: true,
	}

	return res
}
