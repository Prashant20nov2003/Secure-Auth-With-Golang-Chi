package function

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

func LoadRSAPrivateKey() (*rsa.PrivateKey, error){
	privKeyBytes, err := ioutil.ReadFile("./key/private_key.pem")
	if err != nil{
		return nil, err
	}

	block, _ := pem.Decode(privKeyBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY"{
		return nil, fmt.Errorf("failed to decode PEM Block Containing Private Key")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func LoadRSAPublicKey() (*rsa.PublicKey, error){
	pubKeyBytes, err := ioutil.ReadFile("./key/public_key.pem")
	if err != nil{
		return nil, err
	}

	block, _ := pem.Decode(pubKeyBytes)
	if block == nil || block.Type != "RSA PUBLIC KEY"{
		return nil, fmt.Errorf("failed to decode PEM Block Containing Public Key")
	}
	
	return x509.ParsePKCS1PublicKey(block.Bytes)
}