package function

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateEmailToken(emailverify_id string,user_id string) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["emailverify_id"] = emailverify_id
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Minute * 5).Unix()

	jwt_token, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return jwt_token, nil
}

func EmailTokenValidation(email_token string) (jwt.MapClaims, error) {
	// Parse jwt token
	token, err := jwt.Parse(email_token, func(token *jwt.Token) (interface{}, error){
		return publicKey, nil
	})
	if err != nil{
		return nil, fmt.Errorf("couldn't parse token %e",err)
	}

	// Check if jwt token valid
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid{
		if float64(time.Now().Unix()) > claims["exp"].(float64){
			return nil, fmt.Errorf("jwt token expired, please login again")
		}else{
			return claims, nil
		}
	}else{
		return nil, fmt.Errorf("token invalid, please login again")
	}
}