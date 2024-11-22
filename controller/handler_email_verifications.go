package controller

import (
	. "betamart/function"
	"betamart/internal/database"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	// "golang.org/x/crypto/bcrypt"
)

func (apiCfg *ApiConfig) FetchEmailVerification(res http.ResponseWriter, req *http.Request) {
	// Get id from url params
	email_verify_params := chi.URLParam(req, "id")
	email_verify_uuid, err := uuid.Parse(email_verify_params)
	if err != nil {
		RespondWithError(res, 400, fmt.Sprintln("Couldn't parse id format", err))
		return
	}

	// Get the time left
	time_left, err := apiCfg.Query.FetchEmailVerification(req.Context(), email_verify_uuid)
	if err != nil {
		RespondWithError(res, 401, fmt.Sprintln("Couldn't fetch time left", err))
		return
	}

	// If time less equal than 0, then email verification is expired
	if time_left <= 0 {
		RespondWithError(res, 401, fmt.Sprintln("Email Verification Expired"))
		return
	}

	// response
	type Response struct {
		TimeLeft int32 `json:"time_left"`
	}
	response := Response{
		TimeLeft: time_left,
	}
	RespondWithJSON(res, 200, response)
}

func (apiCfg *ApiConfig) VerifyEmailVerification(res http.ResponseWriter, req *http.Request) {
	// Decode Parameter
	type parameters struct {
		VerifCode string `json:"verif_code"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(res, 400, fmt.Sprintln("Error parsing json", err))
		return
	}

	// Get id from url params
	email_verif_params := chi.URLParam(req, "id")
	email_verif_uuid, err := uuid.Parse(email_verif_params)
	if err != nil {
		RespondWithError(res, 400, fmt.Sprintln("Couldn't parse id format", err))
		return
	}

	// Begin Transaction
	tx, err := apiCfg.DB.Begin()
	if err != nil {
		RespondWithError(res, 400, fmt.Sprintln("Failed Making Transaction", err))
		return
	}
	qtx := apiCfg.Query.WithTx(tx)
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Check Email Verification
	email_verif, err := qtx.CheckEmailVerification(req.Context(), database.CheckEmailVerificationParams{
		VerifCode:     params.VerifCode,
		EmailverifyID: email_verif_uuid,
	})
	if err != nil {
		RespondWithError(res, 401, fmt.Sprintln("Couldn't check email verification", err))
		return
	} else if email_verif.ResMessage == "Success" {
		// Used For Verified Email
		if email_verif.UsedFor == "Verify Email" {
			// Update User Email
			user, err := qtx.UpdateUserEmail(req.Context(), database.UpdateUserEmailParams{
				UserID:        email_verif.UserID,
				EmailverifyID: email_verif.EmailverifyID,
			})
			if err != nil {
				RespondWithError(res, 403, fmt.Sprintln("Couldn't update user email", err))
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

			// Commit and response
			type VerifyEmailResponse struct {
				Message      string `json:"message"`
				UsedFor      string `json:"used_for"`
				UserResponse        // Embedded struct
			}
			response := VerifyEmailResponse{
				Message: user.Message,
				UsedFor: email_verif.UsedFor,
				UserResponse: UserResponse{
					Username: user.Username,
					Email:    user.Email,
				},
			}
			tx.Commit()
			RespondWithJSON(res, 200, response)
			return
		} else {
			// Generate Email Token
			email_token, err := GenerateEmailToken(email_verif.EmailverifyID.String(), email_verif.UserID)
			if err != nil {
				RespondWithError(res, 403, fmt.Sprintln("Couldn't generate email token", err))
				return
			}

			expCookie := time.Now().Add(5 * time.Minute)
			http.SetCookie(res, &http.Cookie{
				Name:     "EmailToken",
				Value:    email_token,
				Path:     "/",                     // accessible throughout the site
				HttpOnly: true,                    // prevents JavaScript from accessing the cookie
				Secure:   true,                    // ensures the cookie is sent over HTTPS
				SameSite: http.SameSiteStrictMode, // prevents CSRF attacks
				Expires:  expCookie,               // Set the cookie expiration
			})

			// Commit and response
			type OtherResponse struct {
				Username      string `json:"username"`
				EmailverifyID string `json:"emailverify_id"`
				UsedFor       string `json:"used_for"`
			}
			response := OtherResponse{
				Username:      email_verif.Username,
				EmailverifyID: email_verif.EmailverifyID.String(),
				UsedFor:       email_verif.UsedFor,
			}
			tx.Commit()
			RespondWithJSON(res, 200, response)
			return
		}
	} else {
		// Get Error
		tx.Commit()
		RespondWithError(res, 300, fmt.Sprintln(email_verif.ResMessage))
		return
	}
}

func (apiCfg *ApiConfig) GenerateEmailVerification(res http.ResponseWriter, req *http.Request, user database.User) {
	// Decode Parameter
	type parameters struct {
		UsedFor string `json:"used_for"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(res, 400, fmt.Sprintln("Error parsing json here", err))
		return
	}

	// Begin Transaction
	tx, err := apiCfg.DB.Begin()
	if err != nil {
		RespondWithError(res, 400, fmt.Sprintln("Failed Making Transaction", err))
		return
	}
	qtx := apiCfg.Query.WithTx(tx)
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	// Create Email Verification
	digitCode := DigitGenerator()
	hashedCode := PasswordBcrypt(digitCode)
	emailverify, err := qtx.CreateEmailVerification(req.Context(), database.CreateEmailVerificationParams{
		UserID:    user.UserID,
		VerifCode: string(hashedCode),
		UsedFor:   params.UsedFor,
	})
	if err != nil {
		RespondWithError(res, 400, fmt.Sprintln("Cannot create new email verification", err))
		return
	}

	send_email_res := SendEmailVerification(digitCode, user.Email)
	if send_email_res != "Success" {
		RespondWithError(res, 400, fmt.Sprintln("Cannot send email verification or email is not valid", send_email_res))
		return
	}

	// Commit and response
	tx.Commit()
	RespondWithJSON(res, 200, emailverify)
}

func (apiCfg *ApiConfig) ResendEmailVerificationCode(res http.ResponseWriter, req *http.Request) {
	// Get id from url params
	email_verif_params := chi.URLParam(req, "id")
	email_verif_uuid, err := uuid.Parse(email_verif_params)
	if err != nil {
		RespondWithError(res, 400, fmt.Sprintln("Couldn't parse id format", err))
		return
	}

	// Begin Transaction
	tx, err := apiCfg.DB.Begin()
	if err != nil {
		RespondWithError(res, 400, fmt.Sprintln("Failed Making Transaction", err))
		return
	}
	qtx := apiCfg.Query.WithTx(tx)
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	// Resend Email Verification Code
	digitCode := DigitGenerator()
	hashedCode := PasswordBcrypt(digitCode)
	user, err := qtx.ResendEmailVerificationCode(req.Context(), database.ResendEmailVerificationCodeParams{
		EmailverifyID: email_verif_uuid,
		VerifCode:     string(hashedCode),
	})
	if err != nil {
		RespondWithError(res, 400, fmt.Sprintln("Cannot resend email verification code", err))
		return
	}

	send_email_res := SendEmailVerification(digitCode, user.Email)
	if send_email_res != "Success" {
		RespondWithError(res, 400, fmt.Sprintln("Cannot send email verification or email is not valid", send_email_res))
		return
	}

	// Commit and response
	type Response struct {
		Message       string	`json:"message"`
		EmailverifyID string	`json:"emailverify_id"`
		TimeLeft	int32	`json:"time_left"`
	}
	response := Response{
		Message:       "Verification Code Has Been Resended",
		EmailverifyID: user.EmailverifyID.String(),
		TimeLeft: user.TimeLeft,
	}
	tx.Commit()
	RespondWithJSON(res, 200, response)
}
