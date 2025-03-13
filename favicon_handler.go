package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// FaviconHandler manages favicon operations for projects
type FaviconHandler struct {
	BaseDirectory string // Base directory for storing favicons
}

// NewFaviconHandler creates a new FaviconHandler with the specified base directory
func NewFaviconHandler(baseDir string) *FaviconHandler {
	return &FaviconHandler{
		BaseDirectory: baseDir,
	}
}

// FetchFavicon gets a favicon from a website URL
func (fh *FaviconHandler) FetchFavicon(url string) ([]byte, error) {
	if url == "" {
		return nil, errors.New("URL cannot be empty")
	}

	// Ensure URL has a scheme
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	// Create a custom HTTP client with relaxed TLS settings for testing
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 10,
	}

	// Try direct favicon.ico path first
	parsedURL := strings.TrimSuffix(url, "/")
	faviconURL := fmt.Sprintf("%s/favicon.ico", parsedURL)

	resp, err := client.Get(faviconURL)
	if err == nil && resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		return io.ReadAll(resp.Body)
	}

	// If direct path fails, try to get the web page and parse for favicon
	resp, err = client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch website: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch website, status code: %d", resp.StatusCode)
	}

	// Read the body content
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Extract favicon link from HTML
	bodyStr := string(body)
	re := regexp.MustCompile(`<link[^>]*rel=["'](?:shortcut )?icon["'][^>]*href=["']([^"']+)["'][^>]*>`)
	matches := re.FindStringSubmatch(bodyStr)

	if len(matches) < 2 {
		return nil, fmt.Errorf("no favicon found in the HTML")
	}

	faviconURL = matches[1]
	if !strings.HasPrefix(faviconURL, "http") {
		// Handle relative URLs
		if strings.HasPrefix(faviconURL, "//") {
			faviconURL = "https:" + faviconURL
		} else if strings.HasPrefix(faviconURL, "/") {
			faviconURL = fmt.Sprintf("%s%s", parsedURL, faviconURL)
		} else {
			faviconURL = fmt.Sprintf("%s/%s", parsedURL, faviconURL)
		}
	}

	// Fetch the favicon
	resp, err = client.Get(faviconURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch favicon: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch favicon, status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// GenerateSlug creates a URL-friendly slug from a project name
func GenerateSlug(projectName string) string {
	// Convert to lowercase
	slug := strings.ToLower(projectName)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters
	reg := regexp.MustCompile(`[^a-z0-9-]`)
	slug = reg.ReplaceAllString(slug, "")

	// Remove consecutive hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim leading and trailing hyphens
	slug = strings.Trim(slug, "-")

	return slug
}

// GetFaviconPath generates the proper file path for a project favicon
func (fh *FaviconHandler) GetFaviconPath(projectName string) string {
	slug := GenerateSlug(projectName)
	logosDir := filepath.Join(fh.BaseDirectory, "data", "logos", slug)
	return filepath.Join(logosDir, "favicon.png")
}

// SaveFavicon saves favicon data to the filesystem
func (fh *FaviconHandler) SaveFavicon(projectName string, faviconData []byte) (string, error) {
	if len(faviconData) == 0 {
		return "", errors.New("favicon data cannot be empty")
	}

	slug := GenerateSlug(projectName)
	logosDir := filepath.Join(fh.BaseDirectory, "data", "logos", slug)

	// Create the logos directory if it doesn't exist
	if err := os.MkdirAll(logosDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	faviconPath := filepath.Join(logosDir, "favicon.png")

	// Write the favicon file
	if err := os.WriteFile(faviconPath, faviconData, 0644); err != nil {
		return "", fmt.Errorf("failed to write favicon file: %v", err)
	}

	// Return the relative path to the favicon
	relPath, err := filepath.Rel(fh.BaseDirectory, faviconPath)
	if err != nil {
		return faviconPath, nil // Fall back to absolute path if relative path can't be determined
	}

	return relPath, nil
}

// RemoveFavicon deletes a project's favicon
func (fh *FaviconHandler) RemoveFavicon(projectName string) error {
	faviconPath := fh.GetFaviconPath(projectName)

	if _, err := os.Stat(faviconPath); os.IsNotExist(err) {
		return nil // File doesn't exist, so nothing to remove
	}

	if err := os.Remove(faviconPath); err != nil {
		return fmt.Errorf("failed to remove favicon: %v", err)
	}

	// Try to remove the directory if it's empty
	slug := GenerateSlug(projectName)
	logosDir := filepath.Join(fh.BaseDirectory, "data", "logos", slug)

	// Check if directory is empty
	entries, err := os.ReadDir(logosDir)
	if err == nil && len(entries) == 0 {
		// Directory is empty, try to remove it
		if err := os.Remove(logosDir); err != nil {
			// Non-critical error, just log it
			fmt.Printf("Warning: could not remove empty directory %s: %v\n", logosDir, err)
		}
	}

	return nil
}
