package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

// RouteInfo contiene la información que extraemos de cada línea coincidente.
type RouteInfo struct {
	File    string
	Line    int
	URL     string
	NameURL string
}

func main() {
	// Parseamos los parámetros: ruta de directorio y nombre del archivo CSV de salida.
	dirPath := flag.String("path", "./files_project", "Ruta del directorio a analizar")
	outputCSV := flag.String("out", "output.csv", "Nombre del archivo CSV de salida")
	flag.Parse()

	var routes []RouteInfo

	// Expresión regular para encontrar el patrón @Route("...", name="...")
	re := regexp.MustCompile(`@Route\(\s*"([^"]+)"\s*,\s*name\s*=\s*"([^"]+)"\s*\)`)

	// Recorrer el directorio de forma recursiva.
	err := filepath.Walk(*dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Procesamos solo archivos que no sean directorios, terminen en .php y tengan "Controller.php" al final.
		if !info.IsDir() && filepath.Ext(info.Name()) == ".php" && hasControllerSuffix(info.Name()) {
			fileRoutes, err := processFile(path, re)
			if err != nil {
				log.Printf("Error procesando archivo %s: %v", path, err)
				return nil // Continuamos con los demás archivos
			}
			routes = append(routes, fileRoutes...)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error al recorrer el directorio: %v", err)
	}

	// Generamos el CSV de salida.
	err = writeCSV(*outputCSV, routes)
	if err != nil {
		log.Fatalf("Error al escribir CSV: %v", err)
	}

	fmt.Printf("Análisis completado. Datos exportados a %s\n", *outputCSV)
}

// hasControllerSuffix verifica que el nombre del archivo termine en "Controller.php"
func hasControllerSuffix(filename string) bool {
	return len(filename) >= len("Controller.php") && filename[len(filename)-len("Controller.php"):] == "Controller.php"
}

// processFile abre y procesa un archivo, extrayendo la información relevante línea por línea.
func processFile(path string, re *regexp.Regexp) ([]RouteInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var routes []RouteInfo
	scanner := bufio.NewScanner(file)
	lineNum := 1
	for scanner.Scan() {
		line := scanner.Text()
		// Se busca el patrón en la línea.
		matches := re.FindStringSubmatch(line)
		if len(matches) == 3 {
			// matches[1] es el valor de la URL y matches[2] el name.
			routes = append(routes, RouteInfo{
				File:    filepath.Base(path),
				Line:    lineNum,
				URL:     matches[1],
				NameURL: matches[2],
			})
		}
		lineNum++
	}
	if err := scanner.Err(); err != nil {
		return routes, err
	}
	return routes, nil
}

// writeCSV escribe la información extraída en un archivo CSV.
func writeCSV(filename string, routes []RouteInfo) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	// Escribir la cabecera del CSV.
	header := []string{"file", "n_linea", "url", "name_url"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Escribir los datos.
	for _, r := range routes {
		record := []string{r.File, strconv.Itoa(r.Line), r.URL, r.NameURL}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}
