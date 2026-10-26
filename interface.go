package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func gui() {
	a := app.New()
	w := a.NewWindow("Hello")

	w.SetContent(widget.NewButton(
		"Click me",
		func() { println("Hello!") },
	))

	w.ShowAndRun()
}
