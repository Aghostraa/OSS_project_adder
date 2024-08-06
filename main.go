package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

type Project struct {
	Version     int     `json:"version" yaml:"version"`
	Name        string  `json:"name" yaml:"name"`
	DisplayName string  `json:"displayName" yaml:"display_name"`
	Description string  `json:"description,omitempty" yaml:"description,omitempty"`
	Websites    []URL   `json:"websites,omitempty" yaml:"websites,omitempty"`
	Github      []URL   `json:"github,omitempty" yaml:"github,omitempty"`
	Social      *Social `json:"social,omitempty" yaml:"social,omitempty"`
}

type URL struct {
	Url string `json:"url" yaml:"url"`
}

type Social struct {
	Twitter  []URL `json:"twitter,omitempty" yaml:"twitter,omitempty"`
	Telegram []URL `json:"telegram,omitempty" yaml:"telegram,omitempty"`
	Mirror   []URL `json:"mirror,omitempty" yaml:"mirror,omitempty"`
}

type Response struct {
	Message    string `json:"message"`
	Error      string `json:"error,omitempty"`
	LatestFile string `json:"latestFile,omitempty"`
}

var (
	latestFile string
	mutex      sync.Mutex
)

func main() {
	http.HandleFunc("/createProject", createProjectHandler)
	http.HandleFunc("/getLatestFile", getLatestFileHandler)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createProjectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		setCorsHeaders(w)
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	setCorsHeaders(w)

	var project Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Error decoding JSON: %v", err))
		return
	}

	project.Version = 7

	data, err := yaml.Marshal(&project)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error marshalling YAML: %v", err))
		return
	}

	gitDir := "/Users/ahoura/oss-directory"

	firstChar := strings.ToLower(string(project.Name[0]))
	dirPath := filepath.Join(gitDir, "data/projects", firstChar)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error creating directory: %v", err))
		return
	}

	filePath := filepath.Join(dirPath, fmt.Sprintf("%s.yaml", project.Name))
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error writing file: %v", err))
		return
	}

	// Update the latest file name
	mutex.Lock()
	latestFile = fmt.Sprintf("%s.yaml", project.Name)
	mutex.Unlock()

	log.Printf("Project created: %+v\n", project)
	log.Printf("Latest file created: %s\n", latestFile)

	// Run git commands
	if err := runGitCommand("git", "add", "."); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error adding file to git: %v", err))
		return
	}
	if err := runGitCommand("git", "commit", "-m", "Add new project "+project.Name); err != nil {
		if !strings.Contains(err.Error(), "nothing to commit, working tree clean") {
			writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error committing file to git: %v", err))
			return
		}
	}
	if err := runGitCommand("git", "pull", "origin", "main", "--rebase"); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error pulling changes from git: %v", err))
		return
	}
	if err := runGitCommand("git", "push", "origin", "main"); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error pushing changes to git: %v", err))
		return
	}

	writeSuccessResponse(w, "Project created and changes pushed successfully", latestFile)
}

func getLatestFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	setCorsHeaders(w)

	mutex.Lock()
	currentLatestFile := latestFile
	mutex.Unlock()

	writeSuccessResponse(w, "Latest file retrieved successfully", currentLatestFile)
}

func writeSuccessResponse(w http.ResponseWriter, message string, latestFile string) {
	w.WriteHeader(http.StatusOK)
	response := Response{Message: message, LatestFile: latestFile}
	json.NewEncoder(w).Encode(response)
}

func runGitCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	var out, errBuf strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	cmd.Dir = "/Users/ahoura/oss-directory"
	if err := cmd.Run(); err != nil {
		log.Printf("Error running git command: %v\nOutput: %s\nError: %s\n", err, out.String(), errBuf.String())
		return fmt.Errorf("error running git command: %v\nOutput: %s\nError: %s\n", err, out.String(), errBuf.String())
	}
	log.Printf("Git command output: %s\n", out.String())
	return nil
}

func setCorsHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	log.Println(message)
	w.WriteHeader(statusCode)
	response := Response{Error: message}
	json.NewEncoder(w).Encode(response)
}
