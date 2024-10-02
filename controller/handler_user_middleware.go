package controller

import (
	. "betamart/function"
	"betamart/internal/database"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

type MiddlewareHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *ApiConfig) UserMiddleware(auth MiddlewareHandler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Get the authorization header
		jwt_token, err := req.Cookie("Authorization")
		if err != nil{
			RespondWithError(res, 401, fmt.Sprintln("Cannot get authorization:",err))
			return
		}
		// Check the authorization jwt token
		user_id, err := JwtMiddleware(jwt_token.Value)
		if err != nil {
			RespondWithError(res, 401, fmt.Sprintln("Authorization invalid:",err))
			return
		}

		// Check if user exist or not
		user, err := apiCfg.Query.CheckUserValidationById(req.Context(), user_id)
		if err != nil {
			RespondWithError(res, 400, fmt.Sprintln("User is not exist or something wrong:",err))
			return
		}

		// Turn back the function
		auth(res, req, user)
	}
}

func (apiCfg *ApiConfig) ForgotPassword(auth MiddlewareHandler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Get id from url params
		username := chi.URLParam(req, "username")
		id := chi.URLParam(req, "id")

		// Check User Validation By Username
		user, err := apiCfg.Query.CheckUserValidationByUsername(req.Context(),username)
		if err != nil {
			RespondWithError(res, 400, fmt.Sprintln("User is not exist or something wrong:",err))
			return
		}

		// If format url like this http://localhost:8000/api/forgot_password/{username}/{id} then
		var newReq *http.Request
		if id != ""{
			newReq = req
		} else{ // else format url like this http://localhost:8000/api/forgot_password/{username} then
			// Add 'used_for' to req body
			newReqData := map[string]string{
				"used_for": "Change Password",
			}
			newReqJSON, err := json.Marshal(newReqData)
			if err != nil {
				RespondWithError(res, 500, fmt.Sprintln("Error creating request body:", err))
				return
			}
			// Create a new request with the same method, URL, and header, but a new body
			newReq, err = http.NewRequest(req.Method, req.URL.String(), bytes.NewBuffer(newReqJSON))
			if err != nil {
				RespondWithError(res, 500, fmt.Sprintln("Error creating new request",err))
				return
			}
			// Copy headers from the original request
			newReq.Header = req.Header
		}

		// Turn back the function
		auth(res, newReq, user)
	}
}