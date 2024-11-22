package controller

import (
	. "betamart/function"
	"betamart/internal/database"
	"bytes"
	"io"
	"time"

	// "database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
)

type VerifyLinkResponse struct {
	Message       string `json:"message"`
	EmailVerifyId string `json:"emailverify_id"`
}

type UserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (apiCfg *ApiConfig) RegisterUser(res http.ResponseWriter, req *http.Request) {
	// Decode parameters
	type parameters struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(res, 400, fmt.Sprintln("Error parsing json:", err))
		return
	}

	// Create User
	user, err := apiCfg.Query.RegisterUser(req.Context(), database.RegisterUserParams{
		Username: params.Username,
		Email:    params.Email,
		Password: string(PasswordBcrypt(params.Password)),
	})
	if err != nil {
		RespondWithError(res, 402, fmt.Sprintln("User is already exist or something wrong:", err))
		return
	}

	// Generate Email Verification
	internalReq := &http.Request{
		Method: "POST",
		Body: io.NopCloser(bytes.NewBufferString(`{
			"used_for":"Verify Email"
		}`)),
	}
	internalRes := httptest.NewRecorder()
	apiCfg.GenerateEmailVerification(internalRes, internalReq, user)
	var rowRes map[string]interface{}
	err = json.NewDecoder(internalRes.Body).Decode(&rowRes)
	if err != nil {
		RespondWithError(res, 400, fmt.Sprintln("Failed to decode email verification response:", err))
		return
	}
	if internalRes.Code != http.StatusOK {
		err, _ := rowRes["error"].(string)
		RespondWithError(res, 500, fmt.Sprintln("Couldn't get error from generate email verification:", err))
		return
	}

	// Access the emailverify_id field using type assertion
	emailVerifyId, ok := rowRes["emailverify_id"].(string)
	if !ok {
		RespondWithError(res, 402, fmt.Sprintln("Email field not found in response"))
		return
	}

	verify_link_response := VerifyLinkResponse{
		Message:       "Verify link first",
		EmailVerifyId: emailVerifyId,
	}
	RespondWithJSON(res, 200, verify_link_response)
}

func (apiCfg *ApiConfig) LoginUser(res http.ResponseWriter, req *http.Request) {
	// Decode Parameter
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(res, 400, fmt.Sprintln("Error parsing json:", err))
		return
	}

	// Try to login
	user, err := apiCfg.Query.LoginUser(req.Context(), database.LoginUserParams{
		Username: params.Username,
		Crypt:    params.Password,
	})
	if err != nil {
		RespondWithError(res, 402, fmt.Sprintln("Username doesn't exist or wrong password"))
		return
	}

	if !user.Isemailverified {
		// Generate Email Verification
		internalReq := &http.Request{
			Method: "POST",
			Body: io.NopCloser(bytes.NewBufferString(`{
				"used_for":"Verify Email"
			}`)),
		}
		internalRes := httptest.NewRecorder()
		apiCfg.GenerateEmailVerification(internalRes, internalReq, user)
		var rowRes map[string]interface{}
		err = json.NewDecoder(internalRes.Body).Decode(&rowRes)
		if err != nil {
			RespondWithError(res, 400, fmt.Sprintln("Failed to decode email verification response:", err))
			return
		}
		if internalRes.Code != http.StatusOK {
			err, _ := rowRes["error"].(string)
			RespondWithError(res, 500, fmt.Sprintln("Couldn't get error from generate email verification:", err))
			return
		}

		// Access the emailverify_id field using type assertion
		emailVerifyId, ok := rowRes["emailverify_id"].(string)
		if !ok {
			RespondWithError(res, 402, fmt.Sprintln("Email field not found in response"))
			return
		}

		verify_link_response := VerifyLinkResponse{
			Message:       "Verify link first",
			EmailVerifyId: emailVerifyId,
		}
		RespondWithJSON(res, 200, verify_link_response)
		return
	}

	// Finally, generateJwtToken
	jwt_token, err := GenerateJwtToken(user.UserID)
	if err != nil {
		RespondWithError(res, 500, fmt.Sprintln("Something wrong:", err))
		return
	}

	// Set JWT token as an HTTP-only cookie
	expCookie := time.Now().Add(30 * 24 * time.Hour)
	http.SetCookie(res, &http.Cookie{
		Name:     "Authorization",
		Value:    jwt_token,
		Path:     "/",                     // accessible throughout the site
		HttpOnly: true,                    // prevents JavaScript from accessing the cookie
		Secure:   true,                    // ensures the cookie is sent over HTTPS
		SameSite: http.SameSiteStrictMode, // prevents CSRF attacks
		Expires:  expCookie,               // Set the cookie expiration
	})

	// Respond To User
	login_response := UserResponse{
		Username: user.Username,
		Email:    user.Email,
	}
	RespondWithJSON(res, 200, login_response)
}

func (apiCfg *ApiConfig) GetUsername(res http.ResponseWriter, req *http.Request, user database.User) {
	type UsernameResponse struct {
		Username string `json:"username"`
	}
	username := UsernameResponse{
		Username: user.Username,
	}
	RespondWithJSON(res, 200, username)
}
