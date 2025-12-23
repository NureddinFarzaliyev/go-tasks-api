package main

import (
	"log"

	"github.com/NureddinFarzaliyev/go-tasks-api/internal/database"
)

func main() {
	db, err := database.OpenSQLite("tasks.db")
	if err != nil {
		log.Fatal(err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatal(err)
	}

	cfg := config{
		addr: ":3000",
		db:   db,
	}

	app := &application{config: cfg}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
