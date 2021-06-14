package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//Respond Write to the httpWrite the content of data
func Respond(w http.ResponseWriter, data []byte) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
	w.Write([]byte(data))
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func updateJavaScript(port string, host string) {
	input, err := ioutil.ReadFile("static/js/app.tpl.js")
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "###PORT###") {
			lines[i] = "port=\"" + port + "\""
		}
		if strings.Contains(line, "###HOSR###") {
			lines[i] = "host=\"" + host + "\""
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile("static/js/app.js", []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
