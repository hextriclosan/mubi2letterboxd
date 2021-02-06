package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

func main() {

	consoleMode := flag.Bool("console", false, "command line mode")
	flag.Parse()

	if !*consoleMode {
		handleRequests()
	}

	fmt.Print("Input MUBI userID and press Enter: ")
	var mubiUserId string
	if _, err := fmt.Scanf("%s", &mubiUserId); err == nil {
		if _, err := strconv.ParseUint(mubiUserId, 10, 64); err == nil {
			if csv, err := process(mubiUserId); err == nil {
				if err := ioutil.WriteFile(letterboxdCsvFileName, []byte(csv), 0644); err != nil {
					fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
				}
			} else {
				fmt.Fprintf(os.Stderr, "Error occurred: %s\n", err)
			}
		} else {
			fmt.Fprintf(os.Stderr, "%q is not a valid UserId\n", mubiUserId)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Error reading UserID: %s\n", err)
	}

	fmt.Print("Press Enter to exit")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
