package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"

	_ "github.com/lib/pq"
)

type Message struct {
	ID   int
	Text string
}

// FromJSON decodes a serialized JSON record - User{}
func (m *Message) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(m)
}

// ToJSON encodes a User JSON record
func (m *Message) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(m)
}

// PostgreSQL connection details
var (
	Hostname = "localhost"
	Port     = 5432
	Username = "postgres"
	Password = "postgres"
	Database = "messages_rest"
)

func ConnectPostgres() *sql.DB {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		Hostname, Port, Username, Password, Database)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Println(err)
		return nil
	}

	return db
}

// InsertMessage is for adding a new message to the database
func InsertMessage(m Message) bool {
	db := ConnectPostgres()
	if db == nil {
		log.Println("Cannot connect to PostreSQL!")
		db.Close()
		return false
	}
	defer db.Close()

	if IsMessageValid(m) {
		log.Println("Message", m.Text, "already exist!")
		return false
	}

	stmt, err := db.Prepare("INSERT INTO lists(Text) values($1)")
	if err != nil {
		log.Println("Addmessage:", err)
		return false
	}

	stmt.Exec(m.Text)
	return true
}

// ListAllMessages if for returning all messages from the database table
func ListAllMessages() []Message {
	db := ConnectPostgres()
	if db == nil {
		log.Println("Cannot connect to PostreSQL!")
		db.Close()
		return []Message{}
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM lists \n")
	if err != nil {
		log.Println(err)
		return []Message{}
	}

	all := []Message{}
	var c1 int
	var c2 string

	for rows.Next() {
		err = rows.Scan(&c1, &c2)
		temp := Message{c1, c2}
		all = append(all, temp)
	}

	log.Println("All:", all)
	return all
}

func DeleteMessageID(ID int) bool {
	db := ConnectPostgres()
	if db == nil {
		log.Println("Cannot connect to PostreSQL!")
		db.Close()
		return false
	}
	defer db.Close()

	//check is the message ID exist
	t := FindMessageID(ID)
	if t.ID == 0 {
		log.Println("User", ID, "does not exist")
		return false
	}

	stmt, err := db.Prepare("DELETE FROM lists WHERE ID = $1")
	if err != nil {
		log.Println("DeleteMessage:", err)
		return false
	}

	_, err = stmt.Exec(ID)
	if err != nil {
		log.Println("DeleteMessage:", err)
		return false
	}

	return true
}

func DeleteMessageText(text string) bool {
	db := ConnectPostgres()
	if db == nil {
		log.Println("Cannot connect to PostreSQL!")
		db.Close()
		return false
	}
	defer db.Close()

	//check is the message ID exist
	t := FindMessageText(text)
	if t.Text == "" {
		log.Println("Message", text, "does not exist")
		return false
	}

	stmt, err := db.Prepare("DELETE FROM lists WHERE Text = $1")
	if err != nil {
		log.Println("DeleteMessage:", err)
		return false
	}

	_, err = stmt.Exec(text)
	if err != nil {
		log.Println("DeleteMessage:", err)
		return false
	}

	return true
}

// FindUserID if for returning a message record defined by ID
func FindMessageID(ID int) Message {
	db := ConnectPostgres()
	if db == nil {
		log.Println("Cannot connect to PostreSQL!")
		db.Close()
		return Message{}
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM lists WHERE ID = $1\n", ID)
	if err != nil {
		log.Println("Query:", err)
		return Message{}
	}
	defer rows.Close()

	m := Message{}
	var c1 int
	var c2 string

	for rows.Next() {
		err := rows.Scan(&c1, &c2)
		if err != nil {
			log.Println(err)
			return Message{}
		}
		m = Message{c1, c2}
		log.Println("Found message:", m)
	}
	return m
}

// Same as on top, returns message record by text
func FindMessageText(text string) Message {
	db := ConnectPostgres()
	if db == nil {
		log.Println("Cannot connect to PostreSQL!")
		db.Close()
		return Message{}
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM lists WHERE Text = $1\n", text)
	if err != nil {
		log.Println("Query:", err)
		return Message{}
	}
	defer rows.Close()

	m := Message{}
	var c1 int
	var c2 string

	for rows.Next() {
		err := rows.Scan(&c1, &c2)
		if err != nil {
			log.Println(err)
			return Message{}
		}
		m = Message{c1, c2}
		log.Println("Found message:", m)
	}
	return m
}

func IsMessageValid(m Message) bool {
	db := ConnectPostgres()
	if db == nil {
		log.Println("Cannot connect to PostreSQL!")
		db.Close()
		return false
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM lists WHERE Text = $1 \n", m.Text)
	if err != nil {
		log.Println(err)
		return false
	}

	temp := Message{}
	var c1 int
	var c2 string

	for rows.Next() {
		err = rows.Scan(&c1, &c2)
		if err != nil {
			log.Println(err)
			return false
		}
		temp = Message{c1, c2}
	}
	if m.Text == temp.Text {
		return true
	}
	return false
}
