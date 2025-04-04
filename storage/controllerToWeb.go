package storage

import (
	"database/sql"
	. "github.com/Calevin/go_palantir/parser"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
)

// Plantilla principal para mostrar la tabla de Controllers, con HTMX en la columna "Name URL".
var tmpl = template.Must(template.New("table").Parse(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>Rutas de Controllers</title>
  <!-- Incluir HTMX -->
  <script src="https://unpkg.com/htmx.org@1.9.2"></script>
  <style>
    /* Estilos básicos para el modal */
    #modal {
      display: none;
      position: fixed;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background: rgba(0,0,0,0.5);
      z-index: 1000;
    }
    #modal-content {
      background: white;
      margin: 50px auto;
      padding: 20px;
      width: 80%;
      max-width: 800px;
      border-radius: 5px;
    }
    #modal-close {
      float: right;
      cursor: pointer;
      font-size: 20px;
    }
  </style>
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
      <td>
        <!-- El enlace abre el modal y carga el contenido en el div #modal-body -->
        <a href="#"
           hx-get="/twig?name={{ .NameURL }}"
           hx-target="#modal-body"
           hx-swap="innerHTML"
           onclick="document.getElementById('modal').style.display='block'">
           {{ .NameURL }}
        </a>
      </td>
      <td>{{ .Method }}</td>
    </tr>
    {{ end }}
  </table>
  
  <!-- Modal oculto -->
  <div id="modal">
    <div id="modal-content">
      <span id="modal-close" onclick="document.getElementById('modal').style.display='none'">&times;</span>
      <div id="modal-body">
        <!-- Aquí se cargará la tabla de Twig -->
      </div>
    </div>
  </div>
</body>
</html>
`))

// Template para la tabla de resultados de Twig.
var twigTmpl = template.Must(template.New("twigTable").Parse(`
<h2>Resultados de Twig para "{{ .Name }}"</h2>
<table border="1" cellpadding="5">
  <tr>
    <th>File</th>
    <th>Line</th>
    <th>Path Param</th>
  </tr>
  {{ range .Twigs }}
  <tr>
    <td>{{ .File }}</td>
    <td>{{ .Line }}</td>
    <td>{{ .PathParam }}</td>
  </tr>
  {{ end }}
</table>
`))

// TwigResults estructura para pasar datos al template de Twig.
type TwigResults struct {
	Name  string
	Twigs []TwigPathInfo
}

func RunWeb() {
	// Abrimos la conexión a la base de datos SQLite.
	db, err := sql.Open("sqlite3", "output.db")
	if err != nil {
		log.Fatalf("Error abriendo la base de datos: %v", err)
	}
	defer db.Close()

	// Handler para la raíz: muestra la tabla de controllers.
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

		if err := tmpl.Execute(w, routes); err != nil {
			http.Error(w, "Error renderizando la plantilla", http.StatusInternalServerError)
		}
	})

	// Handler para mostrar la tabla de Twig en base al parámetro NameURL.
	http.HandleFunc("/twig", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "Parámetro 'name' faltante", http.StatusBadRequest)
			return
		}
		// Consultamos la tabla de twig (se asume que la tabla se llama "twig_paths")
		rows, err := db.Query("SELECT file, line, path_param FROM twig_paths WHERE path_param = ?", name)
		if err != nil {
			http.Error(w, "Error consultando twig_paths", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var twigs []TwigPathInfo
		for rows.Next() {
			var twig TwigPathInfo
			if err := rows.Scan(&twig.File, &twig.Line, &twig.PathParam); err != nil {
				http.Error(w, "Error leyendo twig_paths", http.StatusInternalServerError)
				return
			}
			twigs = append(twigs, twig)
		}

		result := TwigResults{
			Name:  name,
			Twigs: twigs,
		}
		if err := twigTmpl.Execute(w, result); err != nil {
			http.Error(w, "Error renderizando la plantilla de twig", http.StatusInternalServerError)
			return
		}
	})

	log.Println("Servidor web escuchando en :8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
