package api

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// This funciton takes an http.FileServer
// if the route matches a static file it's passed to the FileServer
// otherwise the index.html is served to allow client side routing to work
func (a *App) StaticHandler(staticPath string, indexPath string, fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(staticPath, r.URL.Path)

		file, err := os.Stat(path)
		if os.IsNotExist(err) || file.IsDir() {
			// Serve the index file
			http.ServeFile(w, r, filepath.Join(staticPath, indexPath))
			return
		}

		if err != nil {
			log.Printf("Error serving static files: %s\n", err)
			http.Error(w, "Internal Server error", http.StatusInternalServerError)
			return
		}

		// Serve the static file
		fs.ServeHTTP(w, r)
	}
}
