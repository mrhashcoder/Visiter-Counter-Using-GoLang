package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/HouzuoGuo/tiedot/db"
	"github.com/gorilla/mux"
)

type count struct {
	counter int
}

func handleApi(hashDb *db.DB, countDocId int) {
	hashRouter := mux.NewRouter().StrictSlash(true)

	hashRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "welcome to counter api")
		incrementCounter(hashDb, countDocId)
		countCol := hashDb.Use("counterDb")
		readback, err := countCol.Read(countDocId)
		if err != nil {
			panic(err)
		}
		s1 := readback["counter"].(string)
		counter1, _ := strconv.Atoi(s1)
		fmt.Fprintf(w, "Total visits till now is : %d", counter1)
	})
	log.Fatal(http.ListenAndServe(":1234", hashRouter))
}

func incrementCounter(hashDb *db.DB, countDocId int) {
	countCol := hashDb.Use("counterDb")
	readback, err := countCol.Read(countDocId)
	if err != nil {
		panic(err)
	}
	s1 := readback["counter"].(string)
	counter, err := strconv.Atoi(s1)
	counter = counter + 1
	s2 := strconv.Itoa(counter)
	countCol.Update(countDocId, map[string]interface{}{
		"counter": s2,
	})
	readback, err = countCol.Read(countDocId)
	if err != nil {
		panic(err)
	}
}

func main() {
	// set databases
	dbDir := "/database"
	os.RemoveAll(dbDir)
	defer os.RemoveAll(dbDir)

	hashDb, err := db.OpenDB(dbDir)
	if err != nil {
		panic(err)
	}
	err = hashDb.Create("counterDb")
	if err != nil {
		panic(err)
	}

	countCol := hashDb.Use("counterDb")

	countDocId, err := countCol.Insert(map[string]interface{}{
		"counter": "24",
	})
	if err != nil {
		panic(err)
	}

	handleApi(hashDb, countDocId)
}
