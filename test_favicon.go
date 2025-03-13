package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// TestFavicon is a test function for the favicon handler implementation
func TestFavicon() {
	// Create a favicon handler with the test directory
	baseDir := "/Users/ahoura/oss-directory" // Use the same base directory
	handler := NewFaviconHandler(baseDir)

	// Test URL to fetch
	testURL := "https://github.com"
	fmt.Println("Fetching favicon from:", testURL)

	// Test fetching favicon
	favicon, err := handler.FetchFavicon(testURL)
	if err != nil {
		log.Fatalf("Failed to fetch favicon: %v", err)
	}
	fmt.Printf("Successfully fetched favicon (%d bytes)\n", len(favicon))

	// Test project name
	testProject := "Test Project 123"
	fmt.Println("Test project name:", testProject)

	// Test slug generation
	slug := GenerateSlug(testProject)
	fmt.Println("Generated slug:", slug)

	// Test favicon path
	path := handler.GetFaviconPath(testProject)
	fmt.Println("Favicon path:", path)

	// Create directory for testing
	logosDir := filepath.Join(baseDir, "data", "logos", slug)
	os.MkdirAll(logosDir, 0755)

	// Test saving favicon
	savePath, err := handler.SaveFavicon(testProject, favicon)
	if err != nil {
		log.Fatalf("Failed to save favicon: %v", err)
	}
	fmt.Println("Favicon saved to:", savePath)

	// Verify file exists
	if _, err := os.Stat(handler.GetFaviconPath(testProject)); os.IsNotExist(err) {
		log.Fatalf("Favicon file doesn't exist after saving")
	}
	fmt.Println("Verified favicon file exists")

	// Test removing favicon
	err = handler.RemoveFavicon(testProject)
	if err != nil {
		log.Fatalf("Failed to remove favicon: %v", err)
	}
	fmt.Println("Favicon removed successfully")

	// Verify file was removed
	if _, err := os.Stat(handler.GetFaviconPath(testProject)); !os.IsNotExist(err) {
		log.Fatalf("Favicon file still exists after removal")
	}
	fmt.Println("Verified favicon file was removed")

	fmt.Println("All tests passed!")
}

func RunTest() {
	fmt.Println("Running favicon handler tests...")
	TestFavicon()
}
