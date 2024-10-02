package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func StorageProxy() {
	const imageDirectory = `./storage`

	http.HandleFunc("/product_photo/", func(w http.ResponseWriter, r *http.Request) {
		// Create the file path
		filePath := filepath.Join(imageDirectory, r.URL.Path)
		// Serve the file
		for _, ext := range []string{".jpg", ".jpeg", ".png", ".gif"} {
			filePathWithExt := filePath + ext
			if _, err := os.Stat(filePathWithExt); err == nil {
				http.ServeFile(w, r, filePathWithExt)
				return
			}
		}

		http.NotFound(w, r)
	})

	port := ":8080" // or any port of your choice
	fmt.Printf("Serving files from %s on http://localhost%s\n", imageDirectory, port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}