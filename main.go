package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

const (
	ListDir      = 0x0001
	UPLOAD_DIR   = "./uploads"
	TEMPLATE_DIR = "./dist"
)

var templates = make(map[string]*template.Template)

func init() {
	fileInfoArr, err := ioutil.ReadDir(TEMPLATE_DIR)
	check(err)
	var templateName, templatePath string
	for _, fileInfo := range fileInfoArr {
		templateName = fileInfo.Name()
		if ext := path.Ext(templateName); ext != ".html" {
			continue
		}
		templatePath = TEMPLATE_DIR + "/" + templateName
		log.Println("Loading template:", templatePath)
		t := template.Must(template.ParseFiles(templatePath))
		templates[templateName] = t
	}
	fmt.Println(templates)
}

func renderHtml(w http.ResponseWriter, tmpl string, locals map[string]interface{}) {
	err := templates[tmpl].Execute(w, locals)
	check(err)
}

func staticDirHandler(mux *http.ServeMux, prefix string, staticDir string, flags int) {
	mux.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
		file := staticDir + r.URL.Path[len(prefix)-1:]
		fmt.Println("file: ", file)
		if (flags & ListDir) == 0 {
			if exists := isExists(file); !exists {
				http.NotFound(w, r)
				return
			}
		}
		http.ServeFile(w, r, file)
	})
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("tem:", templates)
	if r.Method == "GET" {
		fmt.Println("get / index")
		renderHtml(w, "index.html", nil)

	}

}
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func isExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

func main() {

	mux := http.NewServeMux()
	staticDirHandler(mux, "/", "./dist", 0)
	mux.HandleFunc("/index", IndexHandler)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
