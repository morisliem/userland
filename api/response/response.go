package response

func Response(s string) map[string]string {
	res := map[string]string{}
	res["Message"] = s
	return res
}
