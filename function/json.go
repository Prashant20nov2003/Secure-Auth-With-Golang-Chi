package function

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithError(res http.ResponseWriter, code int, msg string){
	if(code > 499){
		log.Println("Request error", code,"with message:", msg)
	}
	type errorResponse struct{
		Err string `json:"error"`
	}

	RespondWithJSON(res,code,errorResponse{
		Err: msg,
	})
}

func RespondWithJSON(res http.ResponseWriter, code int, payload interface{}){
	data, err := json.Marshal(payload)
	if err != nil {
		log.Println("Failed to get json response:", err)
		res.WriteHeader(500)
		return
	}
	res.Header().Add("Content-Type","application/json")
	res.WriteHeader(code)
	res.Write(data)
}