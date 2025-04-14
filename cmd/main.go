package main

import (
	"context"
	"flag"
	_ "fmt"
	"github.com/Calevin/go_palantir/ent/file"
	"log"
	"os"
	"path/filepath"

	"github.com/Calevin/go_palantir/ent" // Ajusta según tu módulo
	"github.com/Calevin/go_palantir/parser"
	_ "github.com/mattn/go-sqlite3" // Importa el driver SQLite
)

func main() {
	// Parseamos los parámetros: ruta de directorio, db sqlite
	dirPath := flag.String("path", "", "Ruta del directorio a analizar")
	outputSQLite := flag.String("out_sqlite", "", "Nombre del archivo de la db SQLIte")
	flag.Parse()
	// Abrimos la conexión a la base de datos usando Ent.
	ctx := context.Background()
	client, errOpenDb := ent.Open("sqlite3", "file:"+*outputSQLite+"?cache=shared&_fk=1")
	if errOpenDb != nil {
		log.Fatalf("Error abriendo la conexión a SQLite: %v", errOpenDb)
	}
	defer client.Close()

	// Ejecutamos la migración para crear las tablas según el esquema.
	if errCreateSchema := client.Schema.Create(ctx); errCreateSchema != nil {
		log.Fatalf("Error creando el esquema: %v", errCreateSchema)
	}

	// Recorrer el directorio de forma recursiva.
	err := filepath.Walk(*dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
			return err
		}

		if !info.IsDir() && filepath.Ext(info.Name()) == ".php" {
			tokens, errParser := parser.TokenizeFile(path)

			if errParser != nil {
				log.Printf("Error al tokenizar el archivo %s: %v", path, errParser)
				// Se continúa con otros archivos.
				return nil
			}

			// Extraemos el nombre del archivo sin ruta.
			fileName := filepath.Base(path)
			// Buscamos si ya existe un File con ese nombre, si no, lo creamos.
			f, err := client.File.Query().Where(file.NameEQ(fileName)).Only(ctx)
			if err != nil {
				// Si no se encuentra, se crea uno.
				f, err = client.File.
					Create().
					SetName(fileName).
					Save(ctx)
				if err != nil {
					log.Printf("Error creando File para %s: %v", fileName, err)
					return nil
				}
			}

			// Preparamos un slice de creadores para insertar los tokens en bloque.
			var bulkCreates []*ent.TokenCreate
			for _, t := range tokens {
				// Cada "t" es una instancia de ent.Token
				tc := client.Token.Create().
					SetLine(t.Line).
					SetOrder(t.Order).
					SetToken(t.Token).
					SetFile(f) // Aquí se establece la FK hacia la entidad File.
				bulkCreates = append(bulkCreates, tc)
			}

			// Si hay tokens para insertar, se realiza la inserción bulk.
			if len(bulkCreates) > 0 {
				_, err := client.Token.CreateBulk(bulkCreates...).Save(ctx)
				if err != nil {
					log.Printf("Error insertando tokens para el archivo %s: %v", path, err)
				} else {
					log.Printf("Se insertaron %d tokens del archivo %s", len(bulkCreates), path)
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error al recorrer el directorio: %v", err)
	}
}
