package main

import (
	"bytes"
	"fmt"
	"html/template"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
)

func ImageTqt(Name string) {
	existingImageFile, err := os.Open(Name)
	if err != nil {
		log.Fatal(err)
	}

	defer existingImageFile.Close()

	// Calling the generic image.Decode() will tell give us the data
	// and type of image it is as a string. We expect "png"
	_, imageType, err := image.Decode(existingImageFile)
	if err != nil {
		log.Fatal("MAIS", err)
	}
	// fmt.Println(imageData)
	fmt.Println(imageType)

	// We only need this because we already read from the file
	// We have to reset the file pointer back to beginning
	existingImageFile.Seek(0, 0)

	// Since we know it is a png already
	// we can call png.Decode() directly
	loadedImage, err := png.Decode(existingImageFile)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("outimage.png")
	if err != nil {
		// Handle error
	}
	defer f.Close()

	// Encode to `PNG` with `DefaultCompression` level
	// then save to file
	err = png.Encode(f, loadedImage)
	if err != nil {
		// Handle error
	}

}

func VHello(w http.ResponseWriter, r *http.Request) {
	// v := r.FormValue("uploadfile")
	file, _, err := r.FormFile("uploadfile")
	defer file.Close()
	if err != nil {
		log.Fatalf("exec : %s", err)
	}

	buf := bytes.NewBuffer(nil)
	v, err := io.Copy(buf, file)
	if err != nil {
		log.Fatalf("exec : %s", err)
	}

	t := template.New("Welcome")
	t = template.Must(t.ParseFiles("./index.html"))
	err = t.ExecuteTemplate(w, "index.html", "")
	if err != nil {
		log.Fatal(err)
	}
}
func VxHello(w http.ResponseWriter, r *http.Request) {
	t := template.New("Welcome")
	t = template.Must(t.ParseFiles("./new.html"))
	err := t.ExecuteTemplate(w, "new.html", "")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Read image from file that already exists
	http.HandleFunc("/", VHello)
	http.HandleFunc("/upload", VxHello)
	// http.HandleFunc("/upload", VHello)
	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
