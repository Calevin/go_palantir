package storage

import (
	"database/sql"
	. "github.com/Calevin/go_palantir/parser"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
)

// Definimos una plantilla HTML para mostrar los datos en forma de tabla.
var tmpl = template.Must(template.New("table").Parse(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>Rutas de Controllers</title>
</head>
<body>
  <h1>Rutas de Controllers</h1>
  <table border="1" cellpadding="5">
    <tr>
      <th>File</th>
      <th>Line</th>
      <th>URL</th>
      <th>Name URL</th>
      <th>Method</th>
    </tr>
    {{ range . }}
    <tr>
      <td>{{ .File }}</td>
      <td>{{ .Line }}</td>
      <td>{{ .URL }}</td>
      <td>{{ .NameURL }}</td>
      <td>{{ .Method }}</td>
    </tr>
    {{ end }}
  </table>
</body>
</html>
`))

func RunWeb() {
	// Abrimos la conexión a la base de datos SQLite.
	db, err := sql.Open("sqlite3", "output.db")
	if err != nil {
		log.Fatalf("Error abriendo la base de datos: %v", err)
	}
	defer db.Close()

	// Handler para la raíz: consulta los registros y los muestra.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT file, line, url, name_url, method FROM controller_routes")
		if err != nil {
			http.Error(w, "Error consultando la base de datos", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var routes []RouteInfo
		for rows.Next() {
			var route RouteInfo
			if err := rows.Scan(&route.File, &route.Line, &route.URL, &route.NameURL, &route.Method); err != nil {
				http.Error(w, "Error leyendo datos", http.StatusInternalServerError)
				return
			}
			routes = append(routes, route)
		}

		// Renderizamos la plantilla con los datos obtenidos.
		if err := tmpl.Execute(w, routes); err != nil {
			http.Error(w, "Error renderizando la plantilla", http.StatusInternalServerError)
		}
	})

	// Iniciamos el servidor web en el puerto 8888.
	log.Println("Servidor web escuchando en :8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
