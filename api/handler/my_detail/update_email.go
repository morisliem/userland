package mydetail

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"userland/api/helper"
	"userland/api/jwt"
	"userland/api/response"
	"userland/api/validator"
	"userland/store"
)

type UpdateEmailRequest struct {
	Email string `json:"email"`
}

func UpdateUserEmail(userStore store.UserStore, tokenStore store.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request UpdateEmailRequest

		userId, err := helper.AuthenticateUserAccessToken(r, tokenStore)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response.Unautorized_request(err.Error()))
			return
		}

		atJti, err := jwt.GetAtJtiFromAccessToken(r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response.Unautorized_request(err.Error()))
			return
		}

		err = json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response.Bad_request(err.Error()))
			return
		}

		res, err := request.ValidateRequest()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(response.UnproccesableEntity(res))
			return
		}

		newEmail := store.User{
			Email: request.Email,
		}

		rn := helper.GenerateRandomNumber()
		go helper.SendEmailVerCode(newEmail.Email, rn)

		err = tokenStore.SetEmailVerificationCode(userId, rn)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = tokenStore.SetNewEmail(userId, newEmail.Email)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// removing the current access token in redis
		deleted, err := jwt.DeleteATAuth(atJti, tokenStore)
		if err != nil || deleted == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// removing the current refresh token in redis
		itHas, err := tokenStore.HasRefreshToken(atJti)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if itHas {
			deleted, err = jwt.DeleteRTAuth(atJti, tokenStore)
			if err != nil || deleted == 0 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		// Remove the other session as well
		listOfSid, err := userStore.GetSessionsId(r.Context(), userId)
		if err != nil {
			if err == sql.ErrNoRows {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response.Bad_request("unable to get session id"))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for _, v := range listOfSid {
			if v != atJti {
				deleted, err := jwt.DeleteATAuth(v, tokenStore)
				if err != nil || deleted == 0 {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

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

		err = userStore.DeleteCurrentSession(r.Context(), atJti)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
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

func (uer *UpdateEmailRequest) ValidateRequest() (map[string]string, error) {
	res := map[string]string{}
	err := validator.ValidateEmail(uer.Email)
	if err != nil {
		res["email"] = err.Error()
	}

	if len(res) > 0 {
		return res, errors.New("error")
	} else {
		return res, nil
	}
}
