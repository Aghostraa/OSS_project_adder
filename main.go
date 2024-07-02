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

func main() {
	http.HandleFunc("/createProject", createProjectHandler)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createProjectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var project Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding JSON: %v", err), http.StatusBadRequest)
		return
	}

	project.Version = 7

	data, err := yaml.Marshal(&project)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshalling YAML: %v", err), http.StatusInternalServerError)
		return
	}

	gitDir := "/Users/ahoura/oss-directory"

	firstChar := strings.ToLower(string(project.Name[0]))
	dirPath := filepath.Join(gitDir, "data/projects", firstChar)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		http.Error(w, fmt.Sprintf("Error creating directory: %v", err), http.StatusInternalServerError)
		return
	}

	filePath := filepath.Join(dirPath, fmt.Sprintf("%s.yaml", project.Name))
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		http.Error(w, fmt.Sprintf("Error writing file: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Project created: %+v\n", project)

	// Run git commands
	if err := runGitCommand("git", "add", "."); err != nil {
		http.Error(w, fmt.Sprintf("Error adding file to git: %v", err), http.StatusInternalServerError)
		return
	}
	if err := runGitCommand("git", "commit", "-m", "Add new project "+project.Name); err != nil {
		if !strings.Contains(err.Error(), "nothing to commit, working tree clean") {
			http.Error(w, fmt.Sprintf("Error committing file to git: %v", err), http.StatusInternalServerError)
			return
		}
	}
	if err := runGitCommand("git", "pull", "origin", "main", "--rebase"); err != nil {
		http.Error(w, fmt.Sprintf("Error pulling changes from git: %v", err), http.StatusInternalServerError)
		return
	}
	if err := runGitCommand("git", "push", "origin", "main"); err != nil {
		http.Error(w, fmt.Sprintf("Error pushing changes to git: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Project created and changes pushed successfully")
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
