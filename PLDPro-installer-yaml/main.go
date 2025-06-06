package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Nodes    map[int]string `yaml:"nodes"`
	Settings struct {
		InstallPath string `yaml:"installPath"`
		ServerPort  string `yaml:"serverPort"`
		AuthUserEnv string `yaml:"authUserEnv"`
		AuthPassEnv string `yaml:"authPassEnv"`
		Realm       string `yaml:"realm"`
	} `yaml:"settings"`
}

func loadConfig() (*Config, error) {
	f, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(f, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// var nodes []string
var path string
var cfg *Config

func main() {

	var err error
	cfg, err = loadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	// Set path from config
	path = cfg.Settings.InstallPath

	// Start HTTP server
	http.HandleFunc("/", authMiddleware(formHandler))
	http.HandleFunc("/execute", authMiddleware(executeHandler))
	log.Printf("Listening on :%s...", cfg.Settings.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.Settings.ServerPort, nil))
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("form.html"))
	tmpl.Execute(w, nil)
}

func executeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data first
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get selected nodes directly from form values
	selectedNodes := r.Form["nodes"]

	// If "All Nodes" is selected, get all nodes except "Todos"
	if contains(selectedNodes, "0") {
		selectedNodes = make([]string, 0, len(cfg.Nodes)-1)
		for i := 1; i <= 11; i++ {
			if node, ok := cfg.Nodes[i]; ok {
				selectedNodes = append(selectedNodes, node)
			}
		}
	}

	versoesRemovidas := strings.Fields(r.FormValue("versoesRemovidas"))
	setups := strings.Fields(r.FormValue("setups"))
	portas := strings.Fields(r.FormValue("portas"))
	//indices := parseIndices(r.FormValue("indices"))

	//nodesEscolhidos := selectNodes(indices)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming não suportado!", http.StatusInternalServerError)
		return
	}

	//fmt.Fprintf(w, "Nodes escolhidos: %v\n", nodesEscolhidos)
	//fmt.Fprintf(w, "Len Nodes escolhidos: %d\n", len(nodesEscolhidos))
	fmt.Fprintf(w, "Selected nodes: %v\n", selectedNodes)
	fmt.Fprintf(w, "Len Selected nodes: %d\n", len(selectedNodes))
	flusher.Flush()

	sshConfig, err := createSSHConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro na configuração SSH: %v", err), http.StatusInternalServerError)
		return
	}

	numWorkers := runtime.NumCPU()
	jobs := make(chan string, len(selectedNodes))
	results := make(chan string, len(selectedNodes))

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, results, sshConfig, versoesRemovidas, setups, portas)
	}

	for _, host := range selectedNodes {
		jobs <- host
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Fprint(w, result)
		flusher.Flush()
	}

	fmt.Fprintln(w, "Todas as goroutines concluíram.")
	flusher.Flush()
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
