package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Config maps module names to their Git repo URLs
type Config map[string]string

// loadConfig loads module configuration from a JSON file
func loadConfig(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func main() {
	cfg, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// remove leading slash from the path
		mod := r.URL.Path[1:]

		// only handle go-get requests
		if r.URL.Query().Get("go-get") != "1" {
			http.NotFound(w, r)
			return
		}

		// lookup the module in the config
		repo, ok := cfg[mod]
		if !ok {
			http.NotFound(w, r)
			return
		}

		// generate the meta tag for go-import
		meta := fmt.Sprintf(`<meta name="go-import" content="%s git %s">`, r.Host+"/"+mod, repo)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<!doctype html>
<html>
  <head>
    %s
  </head>
  <body>OK</body>
</html>`, meta)
	})

	log.Println("listening on :5247")
	log.Fatal(http.ListenAndServe("127.0.0.1:5247", nil))
}
