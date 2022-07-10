// mubi2letterboxd is a simple command line utility for user data migration from MUBI to letterboxd.
// With the utility, you can create a .csv file suitable for manual import to Letterboxd.
//
// inspired by the reddit entry by jcunews1
// https://www.reddit.com/r/learnjavascript/comments/auwynr/export_mubi_data/ehcx2zf/

package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type MovieRecord struct {
	Id        int    `json:"id"`
	Review    string `json:"body"`
	WatchedAt int64  `json:"updated_at"`
	Rating    int    `json:"overall"`
	Film      struct {
		Title     string `json:"title"`
		Year      int    `json:"year"`
		Directors []struct {
			Name string `json:"name"`
		} `json:"directors"`
	} `json:"film"`
}

const (
	url                   = "https://mubi.com/services/api/ratings"
	perPage               = "1000"
	letterboxdCsvFileName = "letterboxd.csv"
)

func main() {
	fmt.Print("Input MUBI userID and press Enter: ")
	var mubiUserId string
	if _, err := fmt.Scanf("%s", &mubiUserId); err == nil {
		if _, err := strconv.ParseUint(mubiUserId, 10, 64); err == nil {
			if err := process(mubiUserId); err != nil {
				fmt.Fprintf(os.Stderr, "Error occurred: %s\n", err)
			}
		} else {
			fmt.Fprintf(os.Stderr,"%q is not a valid UserId\n", mubiUserId)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Error reading UserID: %s\n", err)
	}

	fmt.Print("Press Enter to exit")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func process(mubiUserId string) error {

	var movieRecords []MovieRecord
	var csvRows [][]string

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	query := req.URL.Query()
	query.Set("user_id", mubiUserId)
	query.Set("per_page", perPage)

	fmt.Printf("Data for UserID %s will be requested from MUBI server\n", mubiUserId)
	for i := 1; ; i++ {
		fmt.Printf("Requesting for chunk #%d... ", i)
		query.Set("page", strconv.Itoa(i))
		req.URL.RawQuery = query.Encode()

		response, err := client.Do(req)
		if err != nil {
			return err
		}

		if response.StatusCode != http.StatusOK {
			return errors.New(fmt.Sprintf("Server returned status code %d", response.StatusCode))
		}

		jsonFile, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(jsonFile, &movieRecords); err != nil {
			return err
		}
		fmt.Printf("downloaded %d records\n", len(movieRecords))

		if len(movieRecords) == 0 {
			break
		}

		for _, item := range movieRecords {
			csvRows = append(csvRows, generateCsvRow(item))
		}
	}

	if len(csvRows) == 0 {
		fmt.Printf("\nNo records found at MUBI server for UserID %s\n", mubiUserId)
		return nil
	}

	outFile, err := os.Create(letterboxdCsvFileName)
	if err != nil {
		return err
	}
	defer func() {
		err := outFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	csvwriter := csv.NewWriter(outFile)
	defer csvwriter.Flush()
	if err := csvwriter.Write([]string{"tmdbID", "Title", "Year", "Directors", "Rating", "WatchedDate", "Review"}); err != nil {
		return err
	}
	if err := csvwriter.WriteAll(csvRows); err != nil {
		return err
	}

	absPath, err := filepath.Abs(outFile.Name())
	if err != nil {
		return err
	}
	fmt.Printf("\n%d records are saved to %q\n", len(csvRows), absPath)

	return outFile.Sync()
}

func generateCsvRow(r MovieRecord) []string {
	idOut := strconv.Itoa(r.Id)
	titleOut := r.Film.Title
	yearOut := strconv.Itoa(r.Film.Year)

	directors := make([]string, len(r.Film.Directors))
	for i, director := range r.Film.Directors {
		directors[i] = director.Name
	}
	directorsOut := strings.Join(directors, ", ")

	ratingOut := strconv.Itoa(r.Rating)
	timeOut := time.Unix(r.WatchedAt, 0).UTC().Format("2006-01-02")
	reviewOut := r.Review

	return []string{idOut, titleOut, yearOut, directorsOut, ratingOut, timeOut, reviewOut}
}
