package main

import (
	"net/http"
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

/*
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
		if strings.Contains(line, "###HOST###") {
			lines[i] = "host=\"" + host + "\""
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile("static/js/app.js", []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
*/
