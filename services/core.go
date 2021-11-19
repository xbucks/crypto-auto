package services

import (
	gosxnotifier "github.com/deckarep/gosx-notifier"
)

func Notify(title string, message string, link string) {
	note := gosxnotifier.NewNotification(message)
	note.Title = title
	note.Sound = gosxnotifier.Default
	note.Link = link
	note.Push()
}
