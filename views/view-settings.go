package views

import (
	"sync"

	"fyne.io/fyne/v2"
	app "fyne.io/fyne/v2/app"
)

func (v *Views) OpenSettings() error {
	v.WaitGroup.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		// Set up base GUI
		thisApp := app.NewWithID("Test Application")
		mainWindow := thisApp.NewWindow("Hidden Main Window")
		mainWindow.Resize(fyne.NewSize(800, 800))
		mainWindow.SetMaster()
		mainWindow.Hide()
		mainWindow.SetCloseIntercept(func() {
			mainWindow.Hide()
		})
		sabWindow := thisApp.NewWindow("SAB Window")
		sabWindow.Resize(fyne.NewSize(640, 480))
		sabWindow.Hide()
		sabWindow.SetCloseIntercept(func() {
			sabWindow.Hide()
		})
		thisApp.Run()

	}(v.WaitGroup)

	return nil
}
