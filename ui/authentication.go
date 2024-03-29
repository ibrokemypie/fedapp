package ui

import (
	"fmt"

	"github.com/ibrokemypie/fedapp/api"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func AppURLWindow(window *widgets.QMainWindow, settings *core.QSettings) {
	// create a regular widget
	// give it a QVBoxLayout
	// and make it the central widget of the window
	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(widgets.NewQVBoxLayout())
	window.SetCentralWidget(widget)

	// // create a line edit
	// // with a custom placeholder text
	// // and add it to the central widgets layout
	input := widgets.NewQLineEdit(nil)
	input.SetPlaceholderText("Instance domain")

	widget.Layout().AddWidget(input)

	// create a button
	// connect the clicked signal
	// and add it to the central widgets layout
	button := widgets.NewQPushButton2("Authenticate", nil)
	button.ConnectClicked(func(bool) {
		authChan := make(chan string)
		go api.Authenticate(input.Text(), authChan, settings)

		authURL := <-authChan
		authURLWindow(window, authURL, authChan)
	})
	input.ConnectReturnPressed(button.Click)
	widget.Layout().AddWidget(button)
	window.Show()
}

func authURLWindow(window *widgets.QMainWindow, authURL string, authChan chan string) {
	// create a regular widget
	// give it a QVBoxLayout
	// and make it the central widget of the window
	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(widgets.NewQVBoxLayout())
	window.SetCentralWidget(widget)

	authURLLabel := widgets.NewQLabel(nil, window.WindowType())
	authURLLabel.SetText("Please open the following URL in your browser and paste the code below:")
	authURLLabel.SetWordWrap(true)
	widget.Layout().AddWidget(authURLLabel)

	// // create a line edit
	// // with a custom placeholder text
	// // and add it to the central widgets layout
	authURLText := widgets.NewQTextEdit(nil)
	authURLText.SetReadOnly(true)
	authURLText.SetText(authURL)
	authURLText.SetMaximumHeight(5 * authURLText.FontMetrics().LineSpacing())
	widget.Layout().AddWidget(authURLText)

	authCodeEntry := widgets.NewQLineEdit(nil)
	authCodeEntry.SetPlaceholderText("Enter code here")
	widget.Layout().AddWidget(authCodeEntry)

	authButton := widgets.NewQPushButton2("Authenticate", nil)
	authButton.ConnectClicked(func(bool) {
		authChan <- authCodeEntry.Text()

		outcome := <-authChan
		if outcome == "success" {
			fmt.Println("successfully authenticated")
		} else {
			fmt.Println("authentication failed")
		}
	})
	authCodeEntry.ConnectReturnPressed(authButton.Click)
	widget.Layout().AddWidget(authButton)

	window.Show()
}
