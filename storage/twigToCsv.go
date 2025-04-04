package storage

import (
	"encoding/csv"
	. "github.com/Calevin/go_palantir/parser"
	"os"
	"strconv"
)

// WriteTwigCSV escribe la información extraída en un archivo CSV.
func WriteTwigCSV(filename string, routes []TwigPathInfo) error {
	f, err := os.Create("twigs_" + filename)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	// Escribir la cabecera del CSV.
	header := []string{"file", "n_linea", "path"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Escribir los datos.
	for _, r := range routes {
		record := []string{
			r.File,
			strconv.Itoa(r.Line),
			r.PathParam,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}
