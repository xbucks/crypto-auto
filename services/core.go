package services

import (
	gosxnotifier "github.com/deckarep/gosx-notifier"
)

func Notify(title string, message string, link string, sound gosxnotifier.Sound) {
	note := gosxnotifier.NewNotification(message)
	note.Title = title
	note.Sound = sound
	// note.Sound = gosxnotifier.Default
	// note.Sound = gosxnotifier.Blow
	note.Link = link
	note.Push()
}
