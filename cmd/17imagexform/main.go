package main

import (
	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
	"gophercises/pkg/transform"
	"gophercises/pkg/util"
	"html/template"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// TODO(oren): clean up file system interactions. this is gross
// TODO(oren): also clean up interactions with command line app (primitive)

var FStore = "/home/oren/src/gophercises/data/17/images"

var templates = template.Must(template.ParseGlob("data/17/*.html"))

func display(w http.ResponseWriter, page string, data interface{}) {
	err := templates.ExecuteTemplate(w, page+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		panic(err)
	}
	mpFile, header, err := r.FormFile("myFile")
	if err != nil {
		log.Panicf("Error retrieving file: %s", err)
	}
	defer mpFile.Close()
	img, fmt, err := image.Decode(mpFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// resize the image for faster primitive transformations
	newImg := resize.Resize(256, 0, img, resize.Lanczos3)

	file, err := os.Create(path.Join(FStore, header.Filename))
	defer file.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch fmt {
	case "jpeg":
		if err := jpeg.Encode(file, newImg, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "png":
		if err := png.Encode(file, newImg); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	table := transform.NewTransformTable(path.Base(file.Name()), 3, 3,
		transform.IncMode(0, 1),
		transform.SingleN(10))
	display(w, "image", table)
}

func handleTransform(w http.ResponseWriter, r *http.Request) {
	tmp := strings.Split(r.URL.Path[1:], "/")
	infile, mode, n := tmp[1], tmp[2], tmp[3]
	infile = path.Join(FStore, infile)
	ext := filepath.Ext(infile)
	outfile := strings.TrimSuffix(infile, filepath.Ext(infile)) + "_" + mode + "_" + n + ext

	if !util.FileExists(outfile) {
		args := []string{
			"-m", mode,
			"-i", infile,
			"-o", outfile,
			"-n", n,
			"-j", "6",
		}
		command := exec.Command("primitive", args...)
		if err := command.Run(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.ServeFile(w, r, outfile)
}

func handleRefine(w http.ResponseWriter, r *http.Request) {
	tmp := strings.Split(r.URL.Path[1:], "/")
	file := tmp[1]
	mode, _ := strconv.Atoi(tmp[2])
	table := transform.NewTransformTable(file, 2, 2,
		transform.SingleMode(mode),
		transform.IncN(10, 80),
		transform.Downloadable())
	display(w, "image", table)

}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		display(w, "upload", nil)
	case "POST":
		uploadFile(w, r)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/upload", handleUpload).
		Methods("GET", "POST")
	r.PathPrefix("/image/").
		Handler(http.StripPrefix("/image/", http.FileServer(http.Dir(FStore))))
	r.HandleFunc("/transform/{file}/{mode}/{n}", handleTransform)
	r.HandleFunc("/refine/{file}/{mode}", handleRefine)
	log.Fatal(http.ListenAndServe(":3000", r))
}
