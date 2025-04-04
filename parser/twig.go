package parser

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
)

// TwigPathInfo contiene la información extraída de cada aparición del patrón en archivos .twig.
type TwigPathInfo struct {
	File      string // Nombre del archivo .twig
	Line      int    // Número de línea donde se encontró el patrón
	PathParam string // Valor extraído entre comillas luego de path(
}

// TwigParser encapsula el regex para procesar archivos .twig.
type TwigParser struct {
	pathRe *regexp.Regexp
}

// NewTwigParser crea una nueva instancia de TwigParser con el regex precompilado.
// La expresión regular utilizada es:
// path\(\s*(['"])([^'"]+)\1
// Donde:
// - (['"]) captura el tipo de comilla usada (' o ").
// - ([^'"]+) captura el contenido hasta encontrar la comilla de cierre, garantizando que se use la misma.
func NewTwigParser() *TwigParser {
	return &TwigParser{
		pathRe: regexp.MustCompile(`path\(\s*(?:'([^']+)'|"([^"]+)")\s*\)`),
	}
}

// ParseFile procesa un archivo .twig y extrae todas las coincidencias del patrón.
// Retorna un slice de TwigPathInfo con el nombre del archivo, número de línea y el parámetro extraído.
func (tp *TwigParser) ParseFile(path string) ([]TwigPathInfo, error) {
	// Verificamos que el archivo tenga la extensión .twig.
	if filepath.Ext(path) != ".twig" {
		return nil, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var results []TwigPathInfo
	scanner := bufio.NewScanner(file)
	lineNum := 1

	for scanner.Scan() {
		line := scanner.Text()

		// Buscamos el patrón path( con comillas simples o dobles.
		// El grupo 1 captura el contenido entre comillas.
		matches := tp.pathRe.FindStringSubmatch(line)
		if len(matches) == 3 {
			results = append(results, TwigPathInfo{
				File:      filepath.Base(path),
				Line:      lineNum,
				PathParam: matches[1],
			})
		}
		lineNum++
	}

	return results, scanner.Err()
}
