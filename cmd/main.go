package main

import (
	"flag"
	"fmt"
	"github.com/Calevin/go_palantir/parser"
	"github.com/Calevin/go_palantir/storage"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// Parseamos los parámetros: ruta de directorio y nombre del archivo CSV de salida.
	dirPath := flag.String("path", "./files_project", "Ruta del directorio a analizar")
	outputCSV := flag.String("out_csv", "", "Nombre del archivo CSV de salida")
	outputSQLite := flag.String("out_sqlite", "", "Nombre del archivo de la db SQLIte")
	flag.Parse()
	var allControllerRows []parser.RouteInfo

	cp := parser.NewControllerParser()

	// Recorrer el directorio de forma recursiva.
	err := filepath.Walk(*dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Procesamos solo archivos que no sean directorios, terminen en .php
		if !info.IsDir() && filepath.Ext(info.Name()) == ".php" {
			controllerRows, err := cp.ParseFile(path)
			if err != nil {
				log.Printf("Error procesando archivo %s: %v", path, err)
				return nil
			}
			allControllerRows = append(allControllerRows, controllerRows...)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error al recorrer el directorio: %v", err)
	}

	if *outputCSV != "" {
		err = storage.WriteCSV(*outputCSV, allControllerRows)
		if err != nil {
			log.Fatalf("Error al escribir CSV: %v", err)
		}

		fmt.Printf("Análisis completado. Datos exportados a %s\n", *outputCSV)
	}

	if *outputSQLite != "" {
		st, err := storage.NewSQLiteControllerStorage(*outputSQLite)
		if err != nil {
			log.Fatalf("Error al escribir DB SQLite: %v", err)
		}

		errSave := st.SaveControllerRoutes(allControllerRows)
		if errSave != nil {
			log.Fatalf("Error al guardar en DB SQLite: %v", errSave)
		}

		fmt.Printf("Análisis completado. Datos exportados a %s\n", *outputSQLite)
	}
}
