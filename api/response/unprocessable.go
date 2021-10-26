package response

type Unprocessable_Entity struct {
	Fields map[string]string `json:"Fields"`
}

func UnproccesableEntity(in map[string]string) Unprocessable_Entity {
	res := Unprocessable_Entity{
		Fields: in,
	}
	return res
}
