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
	Discord  []URL `json:"discord,omitempty" yaml:"discord,omitempty"`
}

type Response struct {
	Message     string   `json:"message"`
	Error       string   `json:"error,omitempty"`
	LatestFile  string   `json:"latestFile,omitempty"`
	StagedFiles []string `json:"stagedFiles,omitempty"`
}

var (
	latestFile  string
	stagedFiles []string
	addedFiles  []string
	mutex       sync.Mutex
)

func main() {
	if err := pullFromUpstream(); err != nil {
		log.Printf("Warning: Failed to pull from upstream: %v", err)
	}

	http.HandleFunc("/createProject", createProjectHandler)
	http.HandleFunc("/getLatestFile", getLatestFileHandler)
	http.HandleFunc("/getCurrentBranch", getCurrentBranchHandler)
	http.HandleFunc("/changeBranch", changeBranchHandler)
	http.HandleFunc("/getAddedFiles", getAddedFilesHandler)
	http.HandleFunc("/getFileContent", getFileContentHandler)
	http.HandleFunc("/getStagedFiles", getStagedFilesHandler)
	http.HandleFunc("/resetFiles", resetFilesHandler)
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
	if _, err := os.Stat(filePath); err == nil {
		writeErrorResponse(w, http.StatusConflict, fmt.Sprintf("File %s already exists", filePath))
		return
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error writing file: %v", err))
		return
	}

	mutex.Lock()
	latestFile = fmt.Sprintf("%s.yaml", project.Name)
	addedFiles = append(addedFiles, latestFile)
	mutex.Unlock()

	// Only stage the changes, don't commit
	if err := stageChanges(); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error staging changes: %v", err))
		return
	}

	writeSuccessResponse(w, "Project created and changes staged", latestFile)
}

func stageChanges() error {
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = "/Users/ahoura/oss-directory"
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error staging changes: %v\nOutput: %s", err, string(output))
	}

	// Get list of staged files
	statusCmd := exec.Command("git", "diff", "--cached", "--name-only")
	statusCmd.Dir = "/Users/ahoura/oss-directory"
	statusOutput, err := statusCmd.Output()
	if err != nil {
		return fmt.Errorf("error getting staged files: %v", err)
	}

	mutex.Lock()
	stagedFiles = strings.Split(strings.TrimSpace(string(statusOutput)), "\n")
	mutex.Unlock()

	return nil
}

func getStagedFilesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	setCorsHeaders(w)

	mutex.Lock()
	response := struct {
		Files []string `json:"files"`
	}{
		Files: stagedFiles,
	}
	mutex.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func getCurrentBranchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	setCorsHeaders(w)

	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = "/Users/ahoura/oss-directory"
	output, err := cmd.Output()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error getting current branch: %v", err))
		return
	}

	currentBranch := strings.TrimSpace(string(output))
	writeSuccessResponse(w, "Current branch retrieved successfully", currentBranch)
}

func changeBranchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	setCorsHeaders(w)

	branchName := r.URL.Query().Get("branch")
	if branchName == "" {
		writeErrorResponse(w, http.StatusBadRequest, "Branch name is required")
		return
	}

	cmd := exec.Command("git", "checkout", branchName)
	cmd.Dir = "/Users/ahoura/oss-directory"
	output, err := cmd.CombinedOutput()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error changing branch: %v\nOutput: %s", err, string(output)))
		return
	}

	writeSuccessResponse(w, fmt.Sprintf("Successfully changed to branch: %s", branchName), "")
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

func getAddedFilesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	setCorsHeaders(w)

	mutex.Lock()
	response := struct {
		Files []string `json:"files"`
	}{
		Files: addedFiles,
	}
	mutex.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func getFileContentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	setCorsHeaders(w)

	filename := r.URL.Query().Get("filename")
	if filename == "" {
		writeErrorResponse(w, http.StatusBadRequest, "Filename is required")
		return
	}

	gitDir := "/Users/ahoura/oss-directory"
	filePath := filepath.Join(gitDir, "data", "projects", string(filename[0]), filename)

	content, err := os.ReadFile(filePath)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error reading file: %v", err))
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(content)
}

func writeSuccessResponse(w http.ResponseWriter, message string, latestFile string) {
	mutex.Lock()
	currentStagedFiles := stagedFiles
	mutex.Unlock()

	w.WriteHeader(http.StatusOK)
	response := Response{
		Message:     message,
		LatestFile:  latestFile,
		StagedFiles: currentStagedFiles,
	}
	json.NewEncoder(w).Encode(response)
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

func pullFromUpstream() error {
	log.Println("Attempting to pull from upstream repository...")

	// First, fetch the latest changes from upstream
	fetchCmd := exec.Command("git", "fetch", "upstream")
	fetchCmd.Dir = "/Users/ahoura/oss-directory"
	fetchOutput, err := fetchCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error fetching from upstream: %v\nOutput: %s", err, string(fetchOutput))
	}
	log.Printf("Fetch from upstream successful. Output: %s", string(fetchOutput))

	// Now, merge the changes into the current branch
	mergeCmd := exec.Command("git", "merge", "upstream/main")
	mergeCmd.Dir = "/Users/ahoura/oss-directory"
	mergeOutput, err := mergeCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error merging upstream changes: %v\nOutput: %s", err, string(mergeOutput))
	}
	log.Printf("Merge from upstream successful. Output: %s", string(mergeOutput))

	log.Println("Successfully pulled and merged changes from upstream repository.")
	return nil
}

func resetAddedFiles() {
	mutex.Lock()
	addedFiles = []string{} // Clear the added files slice
	latestFile = ""         // Reset the latest file
	mutex.Unlock()
}

func resetFilesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	setCorsHeaders(w)
	resetAddedFiles()
	writeSuccessResponse(w, "Files reset successfully", "")
}
