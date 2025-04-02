package storage

import (
	"encoding/csv"
	. "github.com/Calevin/go_palantir/parser"
	"os"
	"strconv"
)

// WriteCSV escribe la información extraída en un archivo CSV.
func WriteCSV(filename string, routes []RouteInfo) error {
	f, err := os.Create("controllers_" + filename)
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
