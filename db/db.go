package db

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/TRileySchwarz/go-database/models"
	"os"
)

// The globally accessible database object
var DataBase *pg.DB

// Initializes the connection to the database
// defaults to port 5432
func InitDatabase() error {
	//Init connection to the database
	ConnectDatabase()

	err := createSchema(DataBase)
	if err != nil {
		return err
	}

	return nil
}

func InitLocalDatabase() error {
	//Init connection to the database
	ConnectLocalDatabase()

	err := createSchema(DataBase)
	if err != nil {
		return err
	}

	return nil
}

func ConnectDatabase() {
	// Update this with your database options ie password and userName
	DataBase = pg.Connect(&pg.Options{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Addr: "db:5432",
	})
}


func ConnectLocalDatabase() {
	// Update this with your database options ie password and userName
	DataBase = pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "my_password",
	})
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
