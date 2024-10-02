package function

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey

func init() {
	var err error
	privateKey, err = LoadRSAPrivateKey()
	if err != nil {
		fmt.Println("Cannot Load RSA Private Key", err)
	}
	publicKey, err = LoadRSAPublicKey()
	if err != nil{
		fmt.Println("Cannot Load RSA Public Key", err)
	}
}

func GenerateJwtToken(user_id string) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()

	jwt_token, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return jwt_token, nil
}

func JwtMiddleware(jwtToken string) (string, error) {
	// Parse jwt token
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error){
		return publicKey, nil
	})
	if err != nil{
		return "", fmt.Errorf("couldn't parse token %e",err)
	}

	// Check if jwt token valid
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid{
		if float64(time.Now().Unix()) > claims["exp"].(float64){
			return "", fmt.Errorf("jwt token expired, please login again")
		}else{
			return claims["user_id"].(string), nil
		}
	}else{
		return "", fmt.Errorf("token invalid, please login again")
	}
}
