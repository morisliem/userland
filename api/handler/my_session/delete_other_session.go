package mysession

import (
	"encoding/json"
	"fmt"
	"net/http"
	"userland/api/helper"
	"userland/api/jwt"
	"userland/api/response"
	"userland/store"
)

// Unable to remove jwt refresh token because i just store the jwt id for access token in the db
func DeleteOtherSession(userStore store.UserStore, tokenStore store.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := helper.AuthenticateUserAccessToken(r, tokenStore)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response.Unautorized_request(err.Error()))
			return
		}

		// getting the jwt id for access token
		atJti, _, err := jwt.GetAtJtiNRtJtiFromAccessToken(r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response.Unautorized_request(err.Error()))
			return
		}

		// getting the list of session id
		listOfSid, err := userStore.GetSessionsId(r.Context(), userId)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// removing all the jwt access token in the redis except the current jwt id
		for _, v := range listOfSid {
			if v != atJti {
				deleted, err := jwt.DeleteATAuth(v, tokenStore)
				if err != nil || deleted == 0 {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Println("1")
					return
				}

				// checking if the session has refresh token id
				itHas, err := tokenStore.HasRefreshToken(v)
				if err != nil {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				if itHas {
					deleted, err = jwt.DeleteRTAuth(v, tokenStore)
					if err != nil || deleted == 0 {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				}
			}
		}

		err = userStore.DeleteOtherSession(r.Context(), userId, atJti)
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
