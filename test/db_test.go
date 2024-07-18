package main

import (
	"testing"
	"testtask/db"
)

func TestDB(t *testing.T) {
	var message = db.Message{ID: 0, Text: "hello everybody"}
	t.Run("check InsertMessage to add new message 'hello everybody'", func(t *testing.T) {
		dbase := db.ConnectPostgres()
		defer dbase.Close()

		got := db.InsertMessage(message)
		want := true

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("сheck listAll", func(t *testing.T) {
		dbase := db.ConnectPostgres()
		defer dbase.Close()

		got := db.ListAllMessages()
		want := []db.Message{
			{
				ID: 1, Text: "hello world",
			},
			{
				ID: 2, Text: "hello everybody",
			},
		}

		if got == nil {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("check deleteMessage to message 'hello everybody'", func(t *testing.T) {
		dbase := db.ConnectPostgres()
		defer dbase.Close()

		got := db.DeleteMessageText("hello everybody")
		want := true

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
