package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

var baseDir string

func init() {
	// Get the base directory of the project
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("Error determining base directory")
		os.Exit(1)
	}
	baseDir = filepath.Dir(filename)
}

func handleBaseRoute(w http.ResponseWriter, r *http.Request) {
	// Construct the absolute path to the templates directory
	templatesDir := filepath.Join(baseDir, "..", "..", "templates")

	// Parse the template using the absolute path
	tmpl, err := template.ParseFiles(filepath.Join(templatesDir, "index.html"))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Execute the template
	tmpl.Execute(w, nil)
}

func handleFormSubmit(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "HTTP Method not allowed", http.StatusMethodNotAllowed)
	}

	name := r.FormValue("name")

	fmt.Fprintf(w, "Hellow, %s! you name has been registered", name)
}

// logging Middleware to log request details
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request received -> Method: %s URL: %s\n\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {
	port := getRandomPort()

	//Routes
	http.Handle("/", loggingMiddleware(http.HandlerFunc(handleBaseRoute)))
	http.Handle("/submit", loggingMiddleware(http.HandlerFunc(handleFormSubmit)))

	fmt.Printf("Server started on Go server started on http://localhost:%d... \n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

}

func getRandomPort() int {
	return 3000 + rand.Intn(6000)
}
