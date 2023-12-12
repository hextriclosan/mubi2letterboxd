package cli

import (
	"fmt"
	"mubi2letterboxd/shared"
	"os"
)

func ProcessCli() error {
	updateStatus("Input MUBI userID and press Enter: ")
	var mubiUserId string

	if _, err := fmt.Scanf("%s", &mubiUserId); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading UserID: %s\n", err)
		return err
	}

	if err := shared.ValidateMubiUserId(mubiUserId); err != nil {
		fmt.Fprint(os.Stderr, err)
		return err
	}

	if err := shared.Process(mubiUserId, shared.LetterboxdCsvFileName, updateStatus); err != nil {
		fmt.Fprintf(os.Stderr, "Error occurred: %s\n", err)
		return err
	}

	return nil
}

func updateStatus(s string) {
	fmt.Printf(s)
}
