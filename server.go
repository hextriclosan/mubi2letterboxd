package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/google/uuid"
)

var (
	globalMap = make(map[string]string)
)

func homePage(w http.ResponseWriter, r *http.Request) {

	type Data struct {
		//DataFields []string
		UserId     string
		Size       int
		DownloadId string
	}

	var data Data

	r.ParseForm()

	data.UserId = r.FormValue("user_id")
	// for key, value := range r.Form {
	// 	data.DataFields = append(data.DataFields, fmt.Sprintf("%s = %v\n", key, value))
	// }

	if data.UserId != "" {
		if csv, err := process(data.UserId); err == nil {
			//w.Write([]byte(csv))
			data.Size = len(csv)

			// if err := ioutil.WriteFile(letterboxdCsvFileName, []byte(csv), 0644); err != nil {
			// 	fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
			// }

			// w.Header().Set("Content-Disposition", "attachment; filename=WHATEVER_YOU_WANT")
			// w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
			// http.ServeFile(w, r, letterboxdCsvFileName)

			data.DownloadId = uuid.New().String()

			globalMap[data.DownloadId] = csv

		} else {
			fmt.Fprintf(os.Stderr, "Error occurred: %s\n", err)
		}
	}

	parsedTemplate, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Println("Error parsing template: ", err)
		return
	}

	err = parsedTemplate.Execute(w, data)
	if err != nil {
		log.Println("Error executing template: ", err)
		return
	}

}

func getFile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("id")

	if id != "" {
		modtime := time.Now()

		// tell the browser the returned content should be downloaded
		w.Header().Add("Content-Disposition", "attachment; filename="+letterboxdCsvFileName)
		w.Header().Set("Content-Type", "text/csv")

		http.ServeContent(w, r, letterboxdCsvFileName, modtime, bytes.NewReader([]byte(globalMap[id])))

		//delete(sessions, "somekey");
		//w.Header().Set("Content-Disposition", "attachment; filename=WHATEVER_YOU_WANT")

	}
}

func handleRequests() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homePage)
	mux.HandleFunc("/get_file", getFile)
	log.Fatal(http.ListenAndServe(":10000", mux))
}
