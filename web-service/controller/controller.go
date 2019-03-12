package controller

import (
	"fmt"
	"github.com/tracechain/fabric-service/fabricSetup"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

type Application struct {
	Fabric *fabricSetup.FabricSetup
}

func renderTemplate(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	lp := filepath.Join("web-service", "templates", "mainlayout.html")
	tp := filepath.Join("web-service/static", "html", name)
	fmt.Println(tp)
	fmt.Println(r.URL)
	// Return a 404 if the template doesn't exist
	info, err := os.Stat(tp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}
	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	if data == nil {
		resultTemplate, err := template.ParseFiles(tp)
		if err != nil {
		fmt.Println(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
		}

		if err := resultTemplate.ExecuteTemplate(w, name, data); err != nil {
		fmt.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		}
		return
	}
	resultTemplate, err := template.ParseFiles(lp)
	if err != nil {
		fmt.Println(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if err := resultTemplate.ExecuteTemplate(w, name, data); err != nil {
		fmt.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}

}
