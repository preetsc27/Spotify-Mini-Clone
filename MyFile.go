package main

import (
	"os"
	"net/http"
	"path/filepath"
	"encoding/json"
)

// making File info struct
type FileInfo struct {
	Name string
	IsDir bool
	Mode os.FileMode
}

const (
	filePrefix = "/music/"
	root = "./music"
)

func main() {

	// creating the server
	http.HandleFunc("/", playerMainFrame)
	http.HandleFunc(filePrefix, File)
	http.ListenAndServe(":1313", nil)

}

func playerMainFrame (w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "./player.html")
}

func File(w http.ResponseWriter, r *http.Request){
	path := filepath.Join(root, r.URL.Path[len(filePrefix):])
	stat, err := os.Stat(path)

	// error check
	if err != nil{
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	
	// is dir check
	if stat.IsDir() {
		serveDir(w, r, path)
		return
	}

	http.ServeFile(w, r, path)
}
func serveDir(w http.ResponseWriter, r *http.Request, path string) {
	defer func() {
		if err, ok := recover().(error); ok {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}()

	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		panic(err)
	}

	files, err := file.Readdir(-1)
	if err != nil {
		panic(err)
	}

	fileInfos:= make([]FileInfo, len(files), len(files))

	for i := range files{
		fileInfos[i].Name = files[i].Name()
		fileInfos[i].IsDir = files[i].IsDir()
		fileInfos[i].Mode = files[i].Mode()
	}

	j := json.NewEncoder(w)

	if err := j.Encode(&fileInfos); err != nil {
		panic(err)
	}
}


