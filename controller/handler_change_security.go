package controller

import (
	. "betamart/function"
	"betamart/internal/database"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) ChangePassword(res http.ResponseWriter, req *http.Request, user database.User){
	// Decode parameters
	type parameters struct{
		Password		string `json:"password"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(res, 400, fmt.Sprintln("Error parsing json", err))
		return
	}

	// Get id from url params
	email_verif_params:= chi.URLParam(req, "id")
	email_verif_uuid, err := uuid.Parse(email_verif_params)
	if err != nil {
		RespondWithError(res, 400, fmt.Sprintln("Couldn't parse id format", err))
		return
	}

	// Validate Email Token
	email_token, err := req.Cookie("EmailToken")
	if err != nil{
		RespondWithError(res, 400, fmt.Sprintln("Couldn't get email token:", err))
		return
	}

	emailverif, err := EmailTokenValidation(email_token.Value)
	if err != nil {
		RespondWithError(res, 400, fmt.Sprintln("Couldn't parse id format", err))
		return
	}

	// Ensure "emailverify_id" and "user_id" exist in the emailverif map
	emailverifID, ok := emailverif["emailverify_id"].(string)
	if !ok || emailverifID == "" {
		RespondWithError(res, 400, "Email verification ID is missing or invalid")
		return
	}
	userID, ok := emailverif["user_id"].(string)
	if !ok || userID == "" {
		RespondWithError(res, 400, "User ID is missing or invalid")
		return
	}

	if email_verif_params == emailverifID && user.UserID == userID {
		user, err := apiCfg.Query.ChangePassword(req.Context(),database.ChangePasswordParams{
			Password: string(PasswordBcrypt(params.Password)),
			EmailverifyID: email_verif_uuid,
		})
		if err != nil {
			RespondWithError(res, 400, fmt.Sprintln("Couldn't change password", err))
			return
		}

		RespondWithJSON(res, 200, user)
	}
}

