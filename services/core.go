package services

import (
	"log"
	"runtime"

	gosxnotifier "github.com/deckarep/gosx-notifier"
	"gopkg.in/toast.v1"
)

func Notify(title string, message string, link string, sound gosxnotifier.Sound) {
	if runtime.GOOS == "windows" {
		notification := toast.Notification{
			AppID:   "crypto.auto",
			Title:   title,
			Message: message,
			Actions: []toast.Action{
				{"protocol", "Open", ""},
				{"protocol", "Cancel", ""},
			},
		}
		err := notification.Push()
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		note := gosxnotifier.NewNotification(message)
		note.Title = title
		note.Sound = sound
		note.Link = link
		note.Push()
	}
}
