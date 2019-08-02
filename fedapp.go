package main

import (
	"os"

	"github.com/therecipe/qt/core"

	"github.com/ibrokemypie/fedapp/ui"

	"github.com/therecipe/qt/widgets"
)

func main() {
	// needs to be called once before you can start using the QWidgets
	app := widgets.NewQApplication(len(os.Args), os.Args)
	app.SetOrganizationName("ibrokemypie")
	app.SetApplicationName("fedapp")
	settings := core.NewQSettings5(nil)

	// create a window
	// with a minimum size of 250*200
	// and sets the title to "Hello Widgets Example"
	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(250, 200)
	window.SetWindowTitle("Instance Authentication")

	if !settings.Contains("access_token") {
		ui.AppURLWindow(window, settings)
	}

	// start the main Qt event loop
	// and block until app.Exit() is called
	// or the window is closed by the user
	app.Exec()
}
