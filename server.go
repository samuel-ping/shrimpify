package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const TMP_FILE_DIRECTORY = "./tmp/uploads"

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "Hello, world!")
}

func shrinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/shrink" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(32 << 20) // 32 MB
	if err != nil {
		http.Error(
			w,
			"The uploaded image is too big. Please choose an image that is less than 32MB in size.",
			http.StatusBadRequest,
		)
		return
	}

	file, h, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "There was a problem with your image.", http.StatusBadRequest)
		return
	}

	imageBuffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(imageBuffer, file); err != nil {
		http.Error(w, "An internal error occurred: "+err.Error(), http.StatusInternalServerError)
		return
	}

	filetype := http.DetectContentType(imageBuffer.Bytes())
	if filetype != "image/jpeg" && filetype != "image/png" {
		http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image", http.StatusBadRequest)
		return
	}

	err = os.MkdirAll(TMP_FILE_DIRECTORY, 0700)
	if err != nil {
		http.Error(w, "An internal error occurred: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpFilePath := fmt.Sprintf("./tmp/uploads/%d%s", time.Now().UnixNano(), h.Filename)
	tmpFile, err := os.Create(tmpFilePath)
	if err != nil {
		http.Error(w, "An internal error occurred: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, file)
	if err != nil {
		http.Error(w, "An internal error occurred: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/shrink", shrinkHandler)

	fmt.Printf("Starting server on port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
