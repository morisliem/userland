package mydetail

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"userland/api/response"
	"userland/store"
)

type SetPictureRequest struct {
	Picture string `json:"picture"`
}

func SetUserPicture(userStore store.UserStore, tokenStore store.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseMultipartForm(10 << 20)

		file, handler, err := r.FormFile("picture")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Bad_request(err.Error()))
			return
		}
		defer file.Close()

		fileData := map[string]string{
			"file name": handler.Filename,
			"file size": fmt.Sprintf("%v", handler.Size),
		}

		// fmt.Println(handler.Header)

		// tempFile, err := ioutil.TempFile("./img/user_profile_image", handler.Filename)
		// if err != nil {
		// 	w.Header().Set("Content-Type", "application/json")
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	json.NewEncoder(w).Encode(response.Bad_request(err.Error()))
		// 	return
		// }

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println("here")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Bad_request(err.Error()))
			return
		}

		// fmt.Println(fileBytes)

		err = ioutil.WriteFile("./img/test.png", fileBytes, 0777)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Bad_request(err.Error()))
			return
		}

		// tempFile.Write(fileBytes)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(fileData)

	}
}
