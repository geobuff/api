package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"

	crypto_rand "crypto/rand"

	"cloud.google.com/go/errorreporting"
	"github.com/geobuff/api/repo"
	"github.com/geobuff/api/src"
	"github.com/geobuff/api/types"
	"github.com/geobuff/api/utils"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var errorClient *errorreporting.Client
var server *src.Server

func main() {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))

	err = loadConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println("successfully loaded .env config")

	environment := os.Getenv("ENV")
	if environment == types.DEV || environment == types.PROD {
		ctx := context.Background()
		errorClient, err = errorreporting.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"), errorreporting.Config{
			OnError: func(err error) {
				log.Printf("Could not log error: %v", err)
			},
		})

		if err != nil {
			log.Fatal(err)
		}
		defer errorClient.Close()
	}

	err = repo.OpenConnection()
	if err != nil {
		panic(err)
	}
	fmt.Println("successfully connected to database")

	err = runMigrations()
	if err != nil {
		panic(err)
	}
	fmt.Println("successfully ran database migrations")

	ts := utils.NewTranslationService()
	es := utils.NewEmailService()
	vs := utils.NewValidationService()
	server = src.NewServer(ts, es, vs)
	fmt.Println("successfully initialized server")

	log.Fatal(server.Start())
}

var loadConfig = func() error {
	return godotenv.Load()
}

var runMigrations = func() error {
	driver, err := postgres.WithInstance(repo.Connection, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
