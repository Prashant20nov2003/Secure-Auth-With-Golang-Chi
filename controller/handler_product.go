package controller

import (
	. "betamart/function"
	"betamart/internal/database"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func (apiCfg *ApiConfig) PostProduct(res http.ResponseWriter, req *http.Request, user database.User){
	// Decode Parameter
	type parameters struct {
		ProductName	string `json:"product_name"`
		Price				int32 `json:"price"`
		Visibility	bool `json:"visibility"`
	}
	// Convert String To Int
	price, err := strconv.Atoi(req.FormValue("price"))
	if err != nil {
		RespondWithError(res, 400, fmt.Sprintf("Invalid price value: %v", err))
		return
	}
	// Convert String To Bool
	visibility, err := strconv.ParseBool(req.FormValue("visibility"))
	if err != nil {
		RespondWithError(res, 400, fmt.Sprintf("Invalid visibility value: %v", err))
		return
	}
	params := parameters{
		ProductName: req.FormValue("product_name"),
		Price:       int32(price), // Convert int to int32
		Visibility:  visibility,
	}

	err = req.ParseMultipartForm(25 << 20) // 25 MB
	if err != nil {
		RespondWithError(res, 400, fmt.Sprintln("Maximum image size is 25 MB", err))
		return
	}

	// Handle file upload
	file, image_header, err := req.FormFile("product_photo")
	if err != nil {
		RespondWithError(res, 400, fmt.Sprintln("Error retrieving file ",err))
		return
	}
	defer file.Close()

	// Check file type
	var fileType string
	if image_header != nil {
		// Check if file type supported or not
		contentType := image_header.Header.Get("Content-Type")
		validTypes := map[string]string{
			"image/jpg": ".jpg",
			"image/jpeg": ".jpeg",
			"image/png":  ".png",
			"image/pneg": ".pneg",
		}
		var valid bool
		fileType, valid = validTypes[contentType]
		if !valid {
			RespondWithError(res, 400, "Only JPG, JPEG, PNG, and PNEG file formats are allowed.")
			return
		}
	}

	// Begin Transaction
	tx, err := apiCfg.DB.Begin()
	if err != nil {
		RespondWithError(res, 400, fmt.Sprintln("Failed Making Transaction", err))
		return
	}
	qtx := apiCfg.Query.WithTx(tx)
	defer func(){
		if err != nil{
			tx.Rollback()
			return
		}
	}()

	// Store the product data
	product, err := qtx.PostProduct(req.Context(), database.PostProductParams{
		UserID:         user.UserID,
		ProductName:    params.ProductName,
		Price:          params.Price,
		Visibility:     params.Visibility,
	})
	if err != nil {
		RespondWithError(res, 400, fmt.Sprintln("Failed to post product", err))
		return
	}

	// Save the file to the server
	out, err := os.Create(fmt.Sprintf("storage/product_photo/%s%s",product.ProductPhoto.String(),fileType))
	if err != nil {
		RespondWithError(res, 400, fmt.Sprintln("Unable to save the file ", err))
		return
	}
	defer out.Close()

	// Copy the uploaded file data to the created file
	_, err = io.Copy(out, file)
	if err != nil {
		RespondWithError(res, 400, fmt.Sprintln("Failed to save file ", err))
		return
	}

	// Commit and Response
	tx.Commit()
	RespondWithJSON(res, 200, fmt.Sprintln("Successfully post product", product))
}

func (apiCfg *ApiConfig) GetProduct(res http.ResponseWriter, req *http.Request, user database.User){
	isPrivate:= req.URL.Query().Get("isPrivate")
	
	var err error
	if isPrivate == "true" { // Get All Product except user product
		var product []database.Product
		product, err = apiCfg.Query.GetMyProduct(req.Context(),user.UserID)
		if err != nil{
			RespondWithError(res, 400, fmt.Sprintln("Cannot get my product data, please try again", err))
			return
		}
		RespondWithJSON(res, 200, product)
	}else{ // Get My Product
		var product []database.GetPublicProductRow
		product, err = apiCfg.Query.GetPublicProduct(req.Context(),user.UserID)
		if err != nil{
			RespondWithError(res, 400, fmt.Sprintln("Cannot get public product data, please try again", err))
			return
		}
		RespondWithJSON(res, 200, product)
	}
}