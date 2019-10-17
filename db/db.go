package db

import (
	"github.com/TRileySchwarz/go-database/models"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"os"
)

// The globally accessible database object
var DataBase *pg.DB

// Initializes the connection to the database
func InitDatabase() error {
	ConnectDatabase("db:" + os.Getenv("DB_PORT"))

	err := createSchema(DataBase)
	if err != nil {
		return err
	}

	return nil
}

// Helper used to connect to a local hosted database vs the containerized version
// Useful for local tests etc.
func InitLocalDatabase() error {
	ConnectDatabase("")

	err := createSchema(DataBase)
	if err != nil {
		return err
	}

	return nil
}

// Use the address if provided, otherwise default
func ConnectDatabase(addr string) {
	dbConf := &pg.Options{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	}

	if addr != "" {
		dbConf.Addr = addr
	}

	DataBase = pg.Connect(dbConf)
}

// Creates the schema the database will be using, ie rows. This could be done via automated migrations
// or it can be done inside of the main go app like this.
func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*models.User)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp: true,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
