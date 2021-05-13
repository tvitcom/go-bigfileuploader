package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	GBYTE            int64 = 1073741824 // 1GB in bytes
	APP_NAME         string
	MAXGB_FILESIZE   string
	DIR_SEPARATOR    string
	HTTP_ENDPOINT    string
	SECRET_LINK      string
	UPLOAD_DIRECTORY string
	TMPDIR           string
	PUBLIC_DIR       string
)

type PageData struct {
	Title      string
	SecretLink string
	Files      []string
	PubDir     string
}

func init() {

	APP_NAME = os.Getenv("APP_NAME")
	if APP_NAME == "" {
		APP_NAME = "Big file Uploader"
	}
	// fmt.Println("INIT GET:APP_NAME:",APP_NAME)

	MAXGB_FILESIZE = os.Getenv("MAXGB_FILESIZE")
	if MAXGB_FILESIZE == "" {
		MAXGB_FILESIZE = "1"
	}
	// fmt.Println("INIT GET:HTTP_ENDPOINT:",HTTP_ENDPOINT)

	DIR_SEPARATOR = os.Getenv("DIR_SEPARATOR")
	if DIR_SEPARATOR == "" {
		DIR_SEPARATOR = "/"
	}

	HTTP_ENDPOINT = os.Getenv("HTTP_ENDPOINT")
	// fmt.Println("INIT GET:HTTP_ENDPOINT:",HTTP_ENDPOINT)
	if HTTP_ENDPOINT == "" {
		HTTP_ENDPOINT = "127.0.0.1:3000"
	}

	SECRET_LINK = os.Getenv("SECRET_LINK")
	// fmt.Println("INIT GET:SECRET_LINK:",SECRET_LINK)
	if SECRET_LINK == "" {
		SECRET_LINK = "/secretu"
	}

	UPLOAD_DIRECTORY = os.Getenv("UPLOAD_DIRECTORY")
	// fmt.Println("INIT GET:UPLOAD_DIRECTORY:",UPLOAD_DIRECTORY)
	if UPLOAD_DIRECTORY == "" {
		UPLOAD_DIRECTORY = "./uploaded"
	}

	TMPDIR = os.Getenv("TMPDIR")
	// fmt.Println("INIT GET:TMPDIR:",TMPDIR)
	if TMPDIR == "" {
		// err := os.Setenv("TMPDIR", "/home/tmp")
		// if err != nil {
		// 	panic(err)
		// }
		TMPDIR = "./tmp"
	}
	PUBLIC_DIR = os.Getenv("PUBLIC_DIR")
	// fmt.Println("INIT GET:PUBLIC_DIR:",PUBLIC_DIR)
	if PUBLIC_DIR == "" {
		PUBLIC_DIR = "./assets/files"
	}
}

func main() {

	fmt.Println("Http server is start on:", "http://"+HTTP_ENDPOINT)
	fmt.Println("SECRET_LINK:", SECRET_LINK)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc(SECRET_LINK, indexHandler)
	http.HandleFunc(SECRET_LINK+"upload", uploadHandler)

	s := &http.Server{
		Addr:           HTTP_ENDPOINT,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("favicon")
	http.Error(w, "Not Found", http.StatusNotFound)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	// http.Error(w, "Bad request", http.StatusBadRequest)
	pd := PageData{
		Title:      APP_NAME,
		SecretLink: "",
	}
	tmpl, _ := template.ParseFiles("templates/index.htmlt")
	tmpl.Execute(w, pd)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	FileIndex := []string{}
	dir, err := os.Open(PUBLIC_DIR)
	if err != nil {
		panic(err)
		return
	}
	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		panic(err)
		return
	}
	for _, fi := range fileInfos {
		FileIndex = append(FileIndex, fi.Name())
	}
	fmt.Println(FileIndex)

	pd := PageData{
		Title:      APP_NAME,
		SecretLink: SECRET_LINK + "upload",
		Files:      FileIndex,
		PubDir:     PUBLIC_DIR,
	}
	tmpl, _ := template.ParseFiles("templates/uploadindex.htmlt")
	tmpl.Execute(w, pd)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		uploadFile(w, r)
	default:
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

}

func uploadFile(w http.ResponseWriter, r *http.Request) {

	// Max filestorage in runtime buffer
	r.ParseMultipartForm(32 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	// fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	// fmt.Printf("File Size: %+v\n", handler.Size)
	// fmt.Printf("MIME Header: %+v\n", handler.Header["Content-Type"][0])

	// Check file size after uploading process
	var max int64
	max, err = strconv.ParseInt(MAXGB_FILESIZE, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
		return
	}
	calcMax := GBYTE * max
	if handler.Size > calcMax {
		http.Error(w, "413 StatusRequestEntityTooLarge\n413 Размер файла превышен\n", http.StatusRequestEntityTooLarge)
		return
	}

	// Fix directory suffixes in .env params (with or without slash)
	DIR_ENDER := "/"
	if strings.HasSuffix(UPLOAD_DIRECTORY, "/") {
		DIR_ENDER = ""
	} else {
		DIR_ENDER = "/"
	}

	// Create file
	dst, err := os.Create(UPLOAD_DIRECTORY + DIR_ENDER + handler.Filename)
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File is successfully uploaded!\nФайл успешно загружен!\n")
}
