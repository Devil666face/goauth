package main

import (
	. "app/cmd"
	"app/database"
	"app/store"
)

func main() {
	dberr := database.Connect()
	if dberr != nil {
		panic(dberr)
	}
	store.SetStore()

	switch SetCli() {
	case MIGRATE:
		err := Migrate()
		if err != nil {
			panic(err)
		}

	case SUPERUSER:
		err := CreateSuperuser()
		if err != nil {
			panic(err)
		}

	case START:
		err := StartApp()
		if err != nil {
			panic(err)
		}
	}
}
