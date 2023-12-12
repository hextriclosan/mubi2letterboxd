package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"mubi2letterboxd/shared"
)

func ProcessGui() {
	application := app.New()
	window := application.NewWindow("User data migration from MUBI.com to letterboxd.com")

	window.SetFixedSize(true)
	window.Resize(fyne.NewSize(800, 500))

	button := widget.NewButton("Create CSV", nil)
	button.Disable()

	statusLabel := widget.NewLabel("")

	entry := widget.NewEntry()
	entry.SetPlaceHolder("MUBI UserID")

	entry.OnChanged = func(s string) {
		err := shared.ValidateMubiUserId(s)
		if err == nil  {
			button.Enable()
			statusLabel.SetText("Ready")
		} else {
			button.Disable()
			statusLabel.SetText(fmt.Sprint(err))
		}
	}

	button.OnTapped = func() {
		mubiUserId := entry.Text

		statusLabel.SetText("")

		statusUpdater := func(s string) {
			statusLabel.SetText(statusLabel.Text + s)
    	}

		fileDialog := dialog.NewFileSave(func(result fyne.URIWriteCloser, err error) {
			if err == nil && result != nil {
				if err := shared.Process(mubiUserId, result.URI().Path(), statusUpdater); err != nil {
					processError := fmt.Errorf("Error occurred: %s\n", err)
					dialog.ShowError(processError, window)
					statusLabel.SetText(fmt.Sprint(processError))
				}
			}
		}, window)

		fileDialog.SetFileName(shared.LetterboxdCsvFileName)
		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
		fileDialog.Show()
	}

	content := container.NewVBox(entry, button, statusLabel)
	window.SetContent(content)

	window.ShowAndRun()
}
