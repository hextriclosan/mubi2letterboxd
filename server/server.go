// mubi2letterboxd is a simple utility for user data migration from MUBI to letterboxd.
// With the utility, you can create a .csv file suitable for manual import to Letterboxd.
//
// inspired by the reddit entry by jcunews1
// https://www.reddit.com/r/learnjavascript/comments/auwynr/export_mubi_data/ehcx2zf/

package main

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/hextriclosan/mubi2letterboxd/shared"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"
)

//go:embed templates/*
var resources embed.FS

var t = template.Must(template.ParseFS(resources, "templates/*"))

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", homePage)
	addr := ":8080"
	fmt.Printf("Server is listening on %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(os.Stderr, "Error occurred: %s\n", err)
		return
	}

	userId := r.FormValue("user_id")
	if userId != "" {
		if csvRows, err := shared.Process(userId, updateStatus); err == nil {
			var csvData bytes.Buffer
			if err := shared.WriteCsv(&csvData, &csvRows); err != nil {
				fmt.Fprintf(os.Stderr, "Error occurred: %s\n", err)
			}

			w.Header().Set("Content-Disposition", "attachment; filename="+shared.LetterboxdCsvFileName)
			w.Header().Set("Content-Type", "text/csv")

			http.ServeContent(w, r, shared.LetterboxdCsvFileName, time.Now(), bytes.NewReader(csvData.Bytes()))

		} else {
			fmt.Fprintf(os.Stderr, "Error occurred: %s\n", err)
		}
	}

	if err := t.ExecuteTemplate(w, "index.html", nil); err != nil {
		log.Println("Error executing template: ", err)
	}

}

func updateStatus(s string) {
	fmt.Printf(s)
}
