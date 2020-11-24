package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
)

const dbfile = "my.db"

var db *bolt.DB

func init() {
	var err error
	db, err = bolt.Open(dbfile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	defer db.Close()

	var err error

	// var key = flag.String("q", "Question domanda", "test")
	// var value = flag.String("a", "Answer risposta", "test")
	// var bucket = flag.String("b", "mybucket", "Bucket to use")
	var address = flag.String("p", ":8080", "TCP port to use")
	flag.Parse()

	mux := http.NewServeMux()

	mux.HandleFunc("/", fullfilment)
	mux.HandleFunc("/insert", insert)

	err = http.ListenAndServe(*address, mux)
	log.Fatal(err)

}

func insert(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "form.html")
	case "POST":

		var err error

		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		bucket := r.FormValue("b")
		question := r.FormValue("q")
		answer := r.FormValue("a")

		err = db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				return err
			}
			b.Put([]byte(question), []byte(answer))
			return err
		})
		if err != nil {
			log.Println(err)
		}

		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucket))
			v := b.Get([]byte(question))
			fmt.Println(string(v))
			return err
		})
	default:
		fmt.Fprintf(w, "Sorry, only POST method is supported.")
	}

}

func fullfilment(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var err error
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		bucket := r.FormValue("b")
		question := r.FormValue("q")
		//answer := r.FormValue("a")

		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucket))
			v := b.Get([]byte(question))
			fmt.Println(string(v))
			return err
		})
		if err != nil {
			log.Println(err)
		}

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

}
