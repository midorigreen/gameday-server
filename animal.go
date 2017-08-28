package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Animal is `t_animal` table mapping struct
type Animal struct {
	ID   int
	Name string
}

func animalHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET", "":
		animalGET(w, r)
	case "POST":
		animalPOST(w, r)
	}
}

func animalGET(w http.ResponseWriter, r *http.Request) {
	fID := r.URL.Query().Get("id")
	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:3306)/animal", cfg.Animal.User, cfg.Animal.Password, cfg.Animal.Host))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sql := "SELECT id, name FROM t_animal"
	if fID != "" {
		sql = fmt.Sprintf("%s WHERE id=%s", sql, fID)
	}
	rows, err := db.Query(sql)
	if err != nil {
		failed(err, w)
		return
	}
	defer rows.Close()

	var animals []Animal
	for rows.Next() {
		var animal Animal
		rows.Scan(&animal.ID, &animal.Name)

		animals = append(animals, animal)
	}

	jsonAnimal, err := json.Marshal(animals)
	if err != nil {
		failed(err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonAnimal)
}

func animalPOST(w http.ResponseWriter, r *http.Request) {
	fName := r.FormValue("name")
	if fName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("required `name` parameter"))
		return
	}
	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:3306)/animal", cfg.Animal.User, cfg.Animal.Password, cfg.Animal.Host))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	insertAnimal, err := db.Prepare("INSERT INTO name values(?)")
	if err != nil {
		failed(err, w)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		failed(err, w)
		return
	}

	result, err := tx.Stmt(insertAnimal).Exec(fName)
}

func failed(err error, w http.ResponseWriter) {
	log.Print(err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Failed"))
}