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

// RouteInfo contiene la información extraída de cada anotación @Route.
type RouteInfo struct {
	File    string
	Line    int
	URL     string
	NameURL string
	Method  string
}

func main() {
	// Parseamos los parámetros: ruta de directorio y nombre del archivo CSV de salida.
	dirPath := flag.String("path", "./files_project", "Ruta del directorio a analizar")
	outputCSV := flag.String("out", "output.csv", "Nombre del archivo CSV de salida")
	flag.Parse()

	var routes []RouteInfo

	// Expresión regular para encontrar el patrón @Route("...", name="...")
	routeRe := regexp.MustCompile(`@Route\(\s*"([^"]+)"\s*,\s*name\s*=\s*"([^"]+)"\s*\)`)
	// Expresión regular para encontrar la declaración de función: public function nombre(
	funcRe := regexp.MustCompile(`public\s+function\s+(\w+)\s*\(`)
	// Expresión regular para encontrar la declaración de clase: class Nombre...
	classRe := regexp.MustCompile(`class\s+(\w+)`)

	// Recorrer el directorio de forma recursiva.
	err := filepath.Walk(*dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Procesamos solo archivos que no sean directorios, terminen en .php y tengan "Controller.php" al final.
		if !info.IsDir() && filepath.Ext(info.Name()) == ".php" && hasControllerSuffix(info.Name()) {
			fileRoutes, err := processFile(path, routeRe, funcRe, classRe)
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

// processFile procesa un archivo línea por línea y asocia las anotaciones @Route
// al contexto correspondiente (clase o función) según se detecte en el archivo.
func processFile(path string, routeRe, funcRe, classRe *regexp.Regexp) ([]RouteInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var routes []RouteInfo
	// pendingRoutes almacena las anotaciones @Route pendientes de asociar al contexto siguiente.
	var pendingRoutes []RouteInfo

	scanner := bufio.NewScanner(file)
	lineNum := 1
	routeClass := false
	for scanner.Scan() {
		line := scanner.Text()

		// Buscamos la anotación @Route y la guardamos en pendingRoutes.
		matches := routeRe.FindStringSubmatch(line)
		if len(matches) == 3 {
			pendingRoutes = append(pendingRoutes, RouteInfo{
				File:    filepath.Base(path),
				Line:    lineNum,
				URL:     matches[1],
				NameURL: matches[2],
			})
		}

		if !routeClass {
			// Si se encuentra la declaración de una clase, asociamos las anotaciones pendientes.
			classMatches := classRe.FindStringSubmatch(line)
			if len(classMatches) == 2 && len(pendingRoutes) > 0 {
				className := classMatches[1]
				for i := range pendingRoutes {
					pendingRoutes[i].Method = className
				}
				routes = append(routes, pendingRoutes...)
				pendingRoutes = pendingRoutes[:0]

				routeClass = true
			}
		}

		// Si se encuentra la declaración de una función, asociamos las anotaciones pendientes.
		funcMatches := funcRe.FindStringSubmatch(line)
		if len(funcMatches) == 2 && len(pendingRoutes) > 0 {
			funcName := funcMatches[1]
			for i := range pendingRoutes {
				pendingRoutes[i].Method = funcName
			}
			routes = append(routes, pendingRoutes...)
			pendingRoutes = pendingRoutes[:0]
		}
		lineNum++
	}
	if err := scanner.Err(); err != nil {
		return routes, err
	}
	// Si quedaron anotaciones pendientes sin encontrar un contexto, se agregan tal cual.
	routes = append(routes, pendingRoutes...)
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
	header := []string{"file", "n_linea", "url", "name_url", "method"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Escribir los datos.
	for _, r := range routes {
		record := []string{
			r.File,
			strconv.Itoa(r.Line),
			r.URL,
			r.NameURL,
			r.Method,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}
