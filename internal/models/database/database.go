package models

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func Get(connStr string) *gorm.DB {

	conn, err := gorm.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error is %e \n Connection string is %s", err, connStr)
	}
	return conn
}
