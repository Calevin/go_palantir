package parser

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
)

// RouteInfo contiene la información extraída de cada anotación @Route.
type RouteInfo struct {
	File    string
	Line    int
	URL     string
	NameURL string
	Method  string
}

// ControllerParser encapsula los regex para procesar archivos de Controller.
type ControllerParser struct {
	routeRe *regexp.Regexp
	funcRe  *regexp.Regexp
	classRe *regexp.Regexp
}

// NewControllerParser crea una nueva instancia de ControllerParser con los regex precompilados.
func NewControllerParser() *ControllerParser {
	return &ControllerParser{
		routeRe: regexp.MustCompile(`@Route\(\s*"([^"]+)"\s*,\s*name\s*=\s*"([^"]+)"\s*\)`),
		funcRe:  regexp.MustCompile(`public\s+function\s+(\w+)\s*\(`),
		classRe: regexp.MustCompile(`class\s+(\w+)`),
	}
}

// hasControllerSuffix verifica que el archivo tenga el sufijo "Controller.php".
func hasControllerSuffix(filename string) bool {
	return len(filename) >= len("Controller.php") && filename[len(filename)-len("Controller.php"):] == "Controller.php"
}

// ParseFile procesa un archivo y extrae las anotaciones @Route, asociándolas al contexto (clase o función)
// Solo procesa archivos PHP que sean de tipo Controller.
func (cp *ControllerParser) ParseFile(path string) ([]RouteInfo, error) {
	// Verificamos si el archivo es de tipo Controller.
	if filepath.Ext(path) != ".php" || !hasControllerSuffix(filepath.Base(path)) {
		return nil, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var routes []RouteInfo
	var pendingRoutes []RouteInfo
	routeClass := false // Para asegurarnos de que la anotación de clase solo se procese una vez.

	scanner := bufio.NewScanner(file)
	lineNum := 1
	for scanner.Scan() {
		line := scanner.Text()

		// Buscamos la anotación @Route y la agregamos a pendingRoutes.
		if matches := cp.routeRe.FindStringSubmatch(line); len(matches) == 3 {
			pendingRoutes = append(pendingRoutes, RouteInfo{
				File:    filepath.Base(path),
				Line:    lineNum,
				URL:     matches[1],
				NameURL: matches[2],
			})
		}

		// Si aún no se ha procesado la anotación de clase, la buscamos.
		if !routeClass {
			if classMatches := cp.classRe.FindStringSubmatch(line); len(classMatches) == 2 && len(pendingRoutes) > 0 {
				className := classMatches[1]
				for i := range pendingRoutes {
					pendingRoutes[i].Method = className
				}
				routes = append(routes, pendingRoutes...)
				pendingRoutes = pendingRoutes[:0]
				routeClass = true
			}
		}

		// Buscamos la declaración de función.
		if funcMatches := cp.funcRe.FindStringSubmatch(line); len(funcMatches) == 2 && len(pendingRoutes) > 0 {
			funcName := funcMatches[1]
			for i := range pendingRoutes {
				pendingRoutes[i].Method = funcName
			}
			routes = append(routes, pendingRoutes...)
			pendingRoutes = pendingRoutes[:0]
		}
		lineNum++
	}
	// Si quedaron anotaciones pendientes, se agregan sin contexto.
	routes = append(routes, pendingRoutes...)
	return routes, scanner.Err()
}
