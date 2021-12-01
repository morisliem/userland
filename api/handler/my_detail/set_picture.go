package mydetail

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"userland/api/helper"
	"userland/api/response"
	"userland/store"
)

type SetPictureRequest struct {
	Picture string `json:"picture"`
}

func SetUserPicture(userStore store.UserStore, tokenStore store.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := helper.AuthenticateUserAccessToken(r, tokenStore)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response.Unautorized_request(err.Error()))
			return
		}

		r.ParseMultipartForm(10 << 20)

		file, _, err := r.FormFile("picture")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer file.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fileName := fmt.Sprintf("%v", userId) + ".png"

		err = ioutil.WriteFile("docs/img/user/"+fileName, fileBytes, 0777)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = userStore.SetUserPicture(r.Context(), userId, fileName)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response.Success())

	}
}
