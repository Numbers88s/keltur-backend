package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	_ "github.com/lib/pq"
)

// Item is each of the item.
type Item struct {
	id       string
	title    string
	vendor   string
	imageURL string
	category string
	price    float32
}

// DB is database connection.
var DB *sql.DB

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/items", itemsIndex)
	http.HandleFunc("/items/show", itemsShow)
	http.HandleFunc("/item/create", itemsCreate)
	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside handler")
	fmt.Fprintf(w, "Hello world from my Go program!")
}

func itemsCreate(w http.ResponseWriter, r *http.Request) {
	var err error
	DB, err := sql.Open("postgres", "user=keltur password=keltur dbname=keltur sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	id := r.FormValue("id")
	title := r.FormValue("title")
	vendor := r.FormValue("vendor")
	imageURL := r.FormValue("imageURL")
	category := r.FormValue("category")
	if id == "" || title == "" || vendor == "" || imageURL == "" || category == "" {
    http.Error(w, http.StatusText(400), 400)
    return
  }
	price, err := strconv.ParseFloat(r.FormValue("price"), 32)
  if err != nil {
    http.Error(w, http.StatusText(400), 400)
    return
  }

	result, err := DB.Exec("INSERT INTO items VALUES($1, $2, $3, $4, $5, $6)", id, title, vendor, imageURL, category, price)
  if err != nil {
    http.Error(w, http.StatusText(500), 500)
    return
  }

	rowsAffected, err := result.RowsAffected()
	if err != nil {
    http.Error(w, http.StatusText(500), 500)
    return
  }

	 fmt.Fprintf(w, "Item %s created successfully (%d row affected)\n", id, rowsAffected)

}

func itemsShow(w http.ResponseWriter, r *http.Request) {
	var err error
	DB, err := sql.Open("postgres", "user=keltur password=keltur dbname=keltur sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	id := r.FormValue("id")
	fmt.Printf("%+v\n", r.Form)
	fmt.Printf("this is your %s", id)
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	row := DB.QueryRow("SELECT * FROM items WHERE id = $1", id)

	item := new(Item)
	databaseError := row.Scan(&item.id, &item.title, &item.vendor, &item.imageURL, &item.category, &item.price)
	if databaseError == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if databaseError != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fmt.Printf("%s, %s, %s, %s, %s, $%.2f\n", item.id, item.title, item.vendor, item.imageURL, item.category, item.price)
	fmt.Fprintf(w, "%s, %s, %s, %s, %s, $%.2f\n", item.id, item.title, item.vendor, item.imageURL, item.category, item.price)

}

func itemsIndex(w http.ResponseWriter, r *http.Request) {
	var err error
	DB, err := sql.Open("postgres", "user=keltur password=keltur dbname=keltur sslmode=disable")
	// db, err := sql.Open("postgres", "postgres://keltur:keltur@localhost/keltur?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	rows, err := DB.Query("SELECT * FROM items")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	items := make([]*Item, 0)
	for rows.Next() {
		item := new(Item)
		err := rows.Scan(&item.id, &item.title, &item.vendor, &item.imageURL, &item.category, &item.price)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	for _, item := range items {
		fmt.Printf("%s, %s, %s, %s, %s, $%.2f\n", item.id, item.title, item.vendor, item.imageURL, item.category, item.price)
		fmt.Fprintf(w, "%s, %s, %s, %s, %s, $%.2f\n", item.id, item.title, item.vendor, item.imageURL, item.category, item.price)
	}
}
