package storage

import (
	"database/sql"
	. "github.com/Calevin/go_palantir/parser"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
)

// TopTwigInfo almacena el nombre y la cantidad de ocurrencias en twig_paths.
type TopTwigInfo struct {
	PathParam string
	Count     int
}

// PageData contiene los datos para renderizar la página: top twigs y rutas de controllers.
type PageData struct {
	TopTwigs    []TopTwigInfo
	Controllers []RouteInfo
}

// Plantilla principal para mostrar la tabla de Controllers, con HTMX en la columna "Name URL".
var tmpl = template.Must(template.New("table").Parse(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>Dashboard</title>
  <script src="https://unpkg.com/htmx.org@1.9.2"></script>
  <style>
    /* Estilos básicos para el modal */
    #modal {
      display: none;
      position: fixed;
      top: 0; left: 0;
      width: 100%; height: 100%;
      background: rgba(0,0,0,0.5);
      z-index: 1000;
    }
    #modal-content {
      background: #fff;
      margin: 50px auto;
      padding: 20px;
      width: 80%;
      max-width: 800px;
      border-radius: 5px;
      position: relative;
    }
    #modal-close {
      position: absolute;
      top: 10px;
      right: 15px;
      cursor: pointer;
      font-size: 20px;
    }
  </style>
</head>
<body>
  <h1>Top 3 Twig Paths</h1>
  <table border="1" cellpadding="5">
    <tr>
      <th>Path Param</th>
      <th>Count</th>
    </tr>
    {{ range .TopTwigs }}
    <tr>
      <td>{{ .PathParam }}</td>
      <td>{{ .Count }}</td>
    </tr>
    {{ end }}
  </table>

  <h1>Rutas de Controllers</h1>
  <table border="1" cellpadding="5">
    <tr>
      <th>File</th>
      <th>Line</th>
      <th>URL</th>
      <th>Name URL</th>
      <th>Method</th>
    </tr>
    {{ range .Controllers }}
    <tr>
      <td>{{ .File }}</td>
      <td>{{ .Line }}</td>
      <td>{{ .URL }}</td>
      <td>
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

  <!-- Modal para mostrar los resultados de Twig -->
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
	db, err := sql.Open("sqlite3", "output.db")
	if err != nil {
		log.Fatalf("Error abriendo la base de datos: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Consultar el top 3 de twig_paths
		twigRows, err := db.Query("SELECT path_param, COUNT(*) as cnt FROM twig_paths WHERE path_param !='' GROUP BY path_param ORDER BY cnt DESC LIMIT 3")
		if err != nil {
			http.Error(w, "Error consultando twig_paths", http.StatusInternalServerError)
			return
		}
		var topTwigs []TopTwigInfo
		for twigRows.Next() {
			var t TopTwigInfo
			if err := twigRows.Scan(&t.PathParam, &t.Count); err != nil {
				http.Error(w, "Error leyendo twig_paths", http.StatusInternalServerError)
				return
			}
			topTwigs = append(topTwigs, t)
		}
		twigRows.Close()

		// Consultar las rutas de controllers
		ctrlRows, err := db.Query("SELECT file, line, url, name_url, method FROM controller_routes")
		if err != nil {
			http.Error(w, "Error consultando controller_routes", http.StatusInternalServerError)
			return
		}
		var controllers []RouteInfo
		for ctrlRows.Next() {
			var route RouteInfo
			if err := ctrlRows.Scan(&route.File, &route.Line, &route.URL, &route.NameURL, &route.Method); err != nil {
				http.Error(w, "Error leyendo controller_routes", http.StatusInternalServerError)
				return
			}
			controllers = append(controllers, route)
		}
		ctrlRows.Close()

		data := PageData{
			TopTwigs:    topTwigs,
			Controllers: controllers,
		}
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Error renderizando la plantilla", http.StatusInternalServerError)
		}
	})

	// Handler para cargar la tabla de Twig al hacer clic en el link
	http.HandleFunc("/twig", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "Parámetro 'name' faltante", http.StatusBadRequest)
			return
		}
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

		twigTmpl := template.Must(template.New("twigTable").Parse(`
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
		type TwigResults struct {
			Name  string
			Twigs []TwigPathInfo
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
