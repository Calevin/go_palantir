package parser

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/Calevin/go_palantir/ent"
)

// tokenizeLine recibe una línea y la parte en tokens.
// La lógica es la siguiente:
//   - Se consideran como separadores los espacios (o tabulaciones).
//   - Además se "aislan" ciertos caracteres como ; y las comillas (simples y dobles),
//     es decir, cuando se encuentra uno de esos caracteres se finaliza el token actual (si hubiera)
//     y se añade el caracter como token independiente.
func tokenizeLine(line string) []string {
	var tokens []string
	var currentToken []rune

	flush := func() {
		if len(currentToken) > 0 {
			tokens = append(tokens, string(currentToken))
			currentToken = currentToken[:0]
		}
	}

	// Definimos los separadores que queremos "aislar" como tokens independientes.
	isSeparator := func(r rune) bool {
		return r == ';' || r == '"' || r == '\''
	}

	// Recorremos la línea carácter a carácter.
	for _, r := range line {
		// Si encontramos espacio o tabulador, finalizamos el token actual.
		if r == ' ' || r == '\t' {
			flush()
		} else if isSeparator(r) {
			// Si el carácter es uno de los separadores, finalizamos el token actual y
			// agregamos el separador como token.
			flush()
			tokens = append(tokens, string(r))
		} else {
			// Sino, acumulamos el carácter en el token actual.
			currentToken = append(currentToken, r)
		}
	}
	// Agregamos el último token si existe.
	flush()
	return tokens
}

// TokenizeFile abre el archivo ubicado en filePath, lo lee línea a línea y genera
// un slice de *ent.Token con la información de cada token (archivo, línea, orden y token).
func TokenizeFile(filePath string) ([]*ent.Token, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tokens []*ent.Token
	scanner := bufio.NewScanner(file)
	lineNumber := 1

	for scanner.Scan() {
		line := scanner.Text()

		// Se tokeniza la línea usando la función definida.
		lineTokens := tokenizeLine(line)
		// Se guarda cada token, asignándole su orden en la línea (comenzando en 1).
		for i, t := range lineTokens {
			ft := &ent.Token{
				File:  filepath.Base(filePath),
				Line:  lineNumber,
				Order: i + 1,
				Token: t,
			}
			tokens = append(tokens, ft)
		}
		lineNumber++
	}
	if err := scanner.Err(); err != nil {
		return tokens, err
	}
	return tokens, nil
}
