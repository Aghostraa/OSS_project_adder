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
    Version     int      `json:"version"`
    Name        string   `json:"name"`
    DisplayName string   `json:"displayName"`
    Description string   `json:"description"`
    Websites    []URL    `json:"websites"`
    Github      []URL    `json:"github"`
    Social      *Social  `json:"social"`
}

type URL struct {
    Url string `json:"url"`
}

type Social struct {
    Twitter  []URL `json:"twitter"`
    Telegram []URL `json:"telegram"`
    Mirror   []URL `json:"mirror"`
}

func main() {
    http.HandleFunc("/createProject", createProjectHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func createProjectHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var project Project
    if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    project.Version = 7

    data, err := yaml.Marshal(&project)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    firstChar := strings.ToLower(string(project.Name[0]))
    dirPath := filepath.Join("/Users/ahoura/oss-directory/data/projects", firstChar)
    if err := os.MkdirAll(dirPath, 0755); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    filePath := filepath.Join(dirPath, fmt.Sprintf("%s.yaml", project.Name))
    if err := os.WriteFile(filePath, data, 0644); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if err := runGitCommand("git", "add", filePath); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if err := runGitCommand("git", "commit", "-m", "Add new project "+project.Name); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if err := runGitCommand("git", "push", "origin", "main"); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "Project created successfully")
}

func runGitCommand(name string, args ...string) error {
    cmd := exec.Command(name, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Dir = "/Users/ahoura/oss-directory"
    return cmd.Run()
}
