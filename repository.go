package main

import (
	"context"
	"github.com/aeramu/mongolib"
)

var db *mongolib.Database

func initDB() error {
	c, err := mongolib.NewSingletonClient(context.Background(), mongodbURL)
	if err != nil {
		return err
	}
	db = mongolib.NewDatabase(c, "saham")
	return nil
}

func saveTicker(t *ticker) error {
	return db.Coll("income_statement").Save(context.Background(), mongolib.NewObjectID(), t)
}
