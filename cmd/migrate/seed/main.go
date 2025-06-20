package main

import (
	"log"

	"github.com/Prasanthi-Peram/pigee-connect/internal/db"
	"github.com/Prasanthi-Peram/pigee-connect/internal/env"
	"github.com/Prasanthi-Peram/pigee-connect/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/socialnetwork?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)

	db.Seed(store, conn)
}