package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ilikeorangutans/phts/db"
	"github.com/ilikeorangutans/phts/model"
	"github.com/ilikeorangutans/phts/storage"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

func main() {
	fmt.Println("# phts repl starting up...")
	db, err := sqlx.Open("postgres", "user=phts_dev dbname=phts_dev sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	backend := &storage.FileBackend{BaseDir: "tmp"}
	backend.Init()

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Migrating database...")
	m, err := migrate.NewWithDatabaseInstance("file://db/migrate", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		log.Println("Database up to date!")
	} else if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Database migrated!")
	}

	exif.RegisterParsers(mknote.All...)

	runRepl(db, backend)
}

func runRepl(dbx *sqlx.DB, backend storage.Backend) {
	keepGoing := true
	reader := bufio.NewReader(os.Stdin)

	wrappedDB := db.WrapDB(dbx)
	collectionRepo := model.NewCollectionRepository(wrappedDB, backend)
	photoRepo := model.NewPhotoRepository(wrappedDB, backend)

	for keepGoing {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		input = strings.TrimSpace(input)

		fmt.Println(input)

		if strings.HasPrefix(input, "quit") {
			keepGoing = false
		} else if strings.HasPrefix(input, "help") {
			fmt.Println("quit - quit repl")
			fmt.Println("help - this message")
			fmt.Println("lc - list recent collections")
			fmt.Println("cc - create collection")
			fmt.Println("dc - delete collection")
			fmt.Println()
		} else if strings.HasPrefix(input, "lc") {
			collections, err := collectionRepo.Recent(10)
			if err != nil {
				log.Fatal(err)
			}

			for _, c := range collections {
				log.Printf("[%d] %s (%s) %d photos", c.ID, c.Name, c.Slug, c.PhotoCount)
			}
		} else if strings.HasPrefix(input, "cc") {
			fmt.Println("name")
			name, err := reader.ReadString('\n')
			name = strings.TrimSpace(name)
			if err != nil {
				log.Fatal(err)
			}
			slug, err := model.SlugFromString(name)
			if err != nil {
				fmt.Printf("%s", err.Error())
				continue
			}

			col := collectionRepo.Create(name, slug)
			col, err = collectionRepo.Save(col)
			if err != nil {
				fmt.Printf("%s", err.Error())
			} else {

				log.Printf("[%d] %s (%s)", col.ID, col.Name, col.Slug)
			}
		} else if strings.HasPrefix(input, "dc") {
			split := strings.Split(input, " ")
			for _, s := range split[1:] {
				if strings.TrimSpace(s) == "" {
					continue
				}
				id, err := strconv.Atoi(strings.TrimSpace(s))
				if err != nil {
					fmt.Println(err.Error())
				} else {
					col, err := collectionRepo.FindByID(int64(id))
					if err != nil {
						fmt.Println(err.Error())
					} else {
						fmt.Printf("Deleting collection [%d] %s", col.ID, col.Name)
						err = collectionRepo.Delete(col)
						if err != nil {
							fmt.Println(err.Error())
						}
					}
				}
			}
		} else if strings.HasPrefix(input, "ap") {
			split := strings.Split(input, " ")
			id, err := strconv.Atoi(split[1])
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			col, err := collectionRepo.FindByID(int64(id))
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			file, err := os.Open(split[2])
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			data, err := ioutil.ReadAll(file)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			photo, err := collectionRepo.AddPhoto(col, filepath.Base(file.Name()), data)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Printf("[%d] %s, %d renditions\n", photo.ID, photo.Filename, photo.RenditionCount)
				for _, rendition := range photo.Renditions {
					fmt.Printf("  [%d] ", rendition.ID)
				}
			}
		} else if strings.HasPrefix(input, "lp") {
			split := strings.Split(input, " ")
			id, err := strconv.Atoi(split[1])
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			col, err := collectionRepo.FindByID(int64(id))
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			photos, _, err := photoRepo.List(col, db.NewPaginator())
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			for _, photo := range photos {
				fmt.Printf("[%d] %s, %d renditions (created %s)\n", photo.ID, photo.Filename, photo.RenditionCount, photo.CreatedAt)
			}
		}
	}

	fmt.Println("Goodbye!")
}
