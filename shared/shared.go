package shared

import (
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

type StatusUpdater func(string)

type movieRecord struct {
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
    LetterboxdCsvFileName = "letterboxd.csv"
)

func ValidateMubiUserId(mubiUserId string) error {
	if len(mubiUserId) == 0 {
		return fmt.Errorf("MUBI User ID is empty")
	}

	if _, err := strconv.ParseUint(mubiUserId, 10, 64); err != nil {
		return fmt.Errorf("%q is not a valid UserId: %s", mubiUserId, err)
	}

	return nil
}

func Process(mubiUserId string, csvFileName string, updateStatus StatusUpdater) error {
	var movieRecords []movieRecord
	var csvRows [][]string

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	query := req.URL.Query()
	query.Set("user_id", mubiUserId)
	query.Set("per_page", perPage)

	updateStatus(fmt.Sprintf("Data for UserID %s will be requested from MUBI server\n", mubiUserId))
	for i := 1; ; i++ {
		updateStatus(fmt.Sprintf("Requesting for chunk #%d... ", i))
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
		updateStatus(fmt.Sprintf("downloaded %d records\n", len(movieRecords)))

		if len(movieRecords) == 0 {
			break
		}

		for _, item := range movieRecords {
			csvRows = append(csvRows, generateCsvRow(item))
		}
	}

	if len(csvRows) == 0 {
		updateStatus(fmt.Sprintf("No records found at MUBI server for UserID %s\n", mubiUserId))
		return nil
	}

	outFile, err := os.Create(csvFileName)
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
	updateStatus(fmt.Sprintf("%d records are saved to %q\n", len(csvRows), absPath))

	return outFile.Sync()
}

func generateCsvRow(r movieRecord) []string {
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
