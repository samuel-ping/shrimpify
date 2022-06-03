package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const TMP_FILE_DIRECTORY = "./tmp/uploads"
const COMPRESSED_QUALITY = 40

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

	processedImageBuffer, err := processImage(imageBuffer.Bytes(), COMPRESSED_QUALITY)
	if err != nil {
		http.Error(w, "An internal error occurred: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/webp")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(processedImageBuffer)))
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"converted-%s\"", h.Filename))
	w.Write(processedImageBuffer)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/shrink", shrinkHandler)

	fmt.Printf("Starting server on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
