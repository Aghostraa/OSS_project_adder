package main

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "gopkg.in/yaml.v2"
)

type Project struct {
    Version     int      `yaml:"version"`
    Name        string   `yaml:"name"`
    DisplayName string   `yaml:"display_name"`
    Description string   `yaml:"description,omitempty"`
    Websites    []URL    `yaml:"websites,omitempty"`
    Github      []URL    `yaml:"github,omitempty"`
    Social      Social   `yaml:"social,omitempty"`
}

type URL struct {
    Url string `yaml:"url"`
}

type Social struct {
    Twitter  []URL `yaml:"twitter,omitempty"`
    Telegram []URL `yaml:"telegram,omitempty"`
    Mirror   []URL `yaml:"mirror,omitempty"`
}

func main() {
    reader := bufio.NewReader(os.Stdin)

    project := Project{Version: 7}

    project.Name = promptUser(reader, "Enter project name (slug): ")
    project.DisplayName = promptUser(reader, "Enter display name: ")
    project.Description = promptOptional(reader, "Enter description (optional, type n to skip): ")
    project.Websites = promptURLs(reader, "Enter website URL (optional, type n to skip): ")
    project.Github = promptURLs(reader, "Enter GitHub URL (optional, type n to skip): ")
    project.Social.Twitter = promptURLs(reader, "Enter Twitter URL (optional, type n to skip): ")
    project.Social.Telegram = promptURLs(reader, "Enter Telegram URL (optional, type n to skip): ")
    project.Social.Mirror = promptURLs(reader, "Enter Mirror URL (optional, type n to skip): ")

    data, err := yaml.Marshal(&project)
    if err != nil {
        fmt.Printf("Error marshalling to YAML: %v\n", err)
        return
    }

    // Determine the first character of the slug
    firstChar := strings.ToLower(string(project.Name[0]))
    // Construct the directory path based on the first character
    dirPath := filepath.Join("/Users/ahoura/oss-directory/data/projects", firstChar)
    // Create the directory if it doesn't exist
    if err := os.MkdirAll(dirPath, 0755); err != nil {
        fmt.Printf("Error creating directory: %v\n", err)
        return
    }
    // Construct the file path
    filePath := filepath.Join(dirPath, fmt.Sprintf("%s.yaml", project.Name))
    
    err = os.WriteFile(filePath, data, 0644)
    if err != nil {
        fmt.Printf("Error writing file: %v\n", err)
        return
    }

    // Execute Git commands
    if err := runGitCommand("git", "add", filePath); err != nil {
        fmt.Printf("Error adding file to git: %v\n", err)
        return
    }
    if err := runGitCommand("git", "commit", "-m", "Add new project "+project.Name); err != nil {
        fmt.Printf("Error committing file to git: %v\n", err)
        return
    }
    if err := runGitCommand("git", "push"); err != nil {
        fmt.Printf("Error pushing changes to git: %v\n", err)
        return
    }

    fmt.Println("Project added and changes pushed to GitHub.")
}

func promptUser(reader *bufio.Reader, prompt string) string {
    fmt.Print(prompt)
    input, _ := reader.ReadString('\n')
    return strings.TrimSpace(input)
}

func promptOptional(reader *bufio.Reader, prompt string) string {
    input := promptUser(reader, prompt)
    if input == "N/O" {
        return ""
    }
    return input
}

func promptURLs(reader *bufio.Reader, prompt string) []URL {
    var urls []URL
    for {
        input := promptUser(reader, prompt)
        if input == "n" {
            break
        }
        urls = append(urls, URL{Url: input})
    }
    return urls
}

func runGitCommand(name string, args ...string) error {
    cmd := exec.Command(name, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Dir = "/Users/ahoura/oss-directory"  // Update this to the root directory of your forked git repository
    return cmd.Run()
}
