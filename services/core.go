package services

import (
	"runtime"

	gosxnotifier "github.com/deckarep/gosx-notifier"
)

func Notify(title string, message string, link string, sound gosxnotifier.Sound) {
	if runtime.GOOS == "windows" {

	} else {
		note := gosxnotifier.NewNotification(message)
		note.Title = title
		note.Sound = sound
		note.Link = link
		note.Push()
	}
}
