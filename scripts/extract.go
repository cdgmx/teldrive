//go:build ignore

package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	maxRetries = 3
	retryDelay = 2 * time.Second
)

func main() {
	urlFlag := flag.String("url", "", "URL to download from")
	outputFlag := flag.String("output", "", "Output directory")
	offlineFlag := flag.Bool("offline", false, "Use cached assets only (offline mode)")
	validateFlag := flag.Bool("validate", true, "Validate downloaded assets")
	flag.Parse()

	if *outputFlag == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Check for offline mode or cached assets
	if *offlineFlag || (*urlFlag == "" && hasLocalAssets(*outputFlag)) {
		fmt.Println("Using cached UI assets (offline mode)")
		if validateCachedAssets(*outputFlag) {
			fmt.Println("✅ Cached UI assets validated successfully!")
			return
		} else {
			fmt.Println("⚠️  Cached UI assets validation failed")
			if *offlineFlag {
				fmt.Println("❌ Offline mode requested but no valid cached assets found")
				os.Exit(1)
			}
		}
	}

	if *urlFlag == "" {
		fmt.Println("❌ No URL provided and no valid cached assets found")
		flag.Usage()
		os.Exit(1)
	}

	// Clean and create output directory
	if err := os.RemoveAll(*outputFlag); err != nil {
		fmt.Printf("Error removing directory: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(*outputFlag, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}

	// Download with retry logic
	fmt.Printf("Downloading UI assets from %s\n", *urlFlag)
	if err := downloadWithRetry(*urlFlag, *outputFlag, maxRetries); err != nil {
		fmt.Printf("❌ Download failed after %d retries: %v\n", maxRetries, err)
		
		// Try to use any existing backup/cached assets
		if tryFallbackAssets(*outputFlag) {
			fmt.Println("✅ Using fallback UI assets")
			return
		}
		
		fmt.Println("❌ No fallback assets available")
		os.Exit(1)
	}

	// Validate downloaded assets
	if *validateFlag && !validateAssets(*outputFlag) {
		fmt.Println("⚠️  Downloaded assets validation failed, but continuing...")
	}

	fmt.Println("✅ UI assets extracted successfully!")
}

func downloadWithRetry(url, outputDir string, maxRetries int) error {
	var lastErr error
	
	for attempt := 1; attempt <= maxRetries; attempt++ {
		fmt.Printf("Attempt %d/%d...\n", attempt, maxRetries)
		
		err := downloadAndExtract(url, outputDir)
		if err == nil {
			return nil
		}
		
		lastErr = err
		fmt.Printf("Attempt %d failed: %v\n", attempt, err)
		
		if attempt < maxRetries {
			fmt.Printf("Retrying in %v...\n", retryDelay)
			time.Sleep(retryDelay)
		}
	}
	
	return fmt.Errorf("all %d attempts failed, last error: %w", maxRetries, lastErr)
}

func downloadAndExtract(url, outputDir string) error {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	// Add GitHub token if available (for rate limiting)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "token "+token)
		fmt.Println("Using GitHub token for authentication")
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("downloading: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "ui-download-*.zip")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Download with progress (for large files)
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return fmt.Errorf("saving download: %w", err)
	}

	// Extract zip file
	return extractZip(tmpFile.Name(), outputDir)
}

func extractZip(zipPath, outputDir string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("opening zip: %w", err)
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(outputDir, file.Name)

		// Security check: prevent zip slip
		if !filepath.HasPrefix(path, filepath.Clean(outputDir)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return fmt.Errorf("creating directory for file: %w", err)
		}

		dstFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return fmt.Errorf("creating file: %w", err)
		}

		srcFile, err := file.Open()
		if err != nil {
			dstFile.Close()
			return fmt.Errorf("opening zip entry: %w", err)
		}

		_, err = io.Copy(dstFile, srcFile)
		dstFile.Close()
		srcFile.Close()
		if err != nil {
			return fmt.Errorf("extracting file: %w", err)
		}
	}

	return nil
}

func hasLocalAssets(outputDir string) bool {
	// Check if ui/dist directory exists and has content
	if info, err := os.Stat(outputDir); err == nil && info.IsDir() {
		if entries, err := os.ReadDir(outputDir); err == nil && len(entries) > 0 {
			return true
		}
	}
	return false
}

func validateCachedAssets(outputDir string) bool {
	return validateAssets(outputDir)
}

func validateAssets(outputDir string) bool {
	// Basic validation: check for essential files
	essentialFiles := []string{"index.html"}
	
	for _, file := range essentialFiles {
		path := filepath.Join(outputDir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("Missing essential file: %s\n", file)
			return false
		}
	}
	
	return true
}

func tryFallbackAssets(outputDir string) bool {
	// Try to find any cached assets in common locations
	fallbackPaths := []string{
		"ui/dist.backup",
		"ui/dist.cache",
		".ui-cache",
	}
	
	for _, fallbackPath := range fallbackPaths {
		if hasLocalAssets(fallbackPath) {
			fmt.Printf("Found fallback assets in %s\n", fallbackPath)
			// Copy fallback assets to output directory
			return copyDir(fallbackPath, outputDir)
		}
	}
	
	return false
}

func copyDir(src, dst string) bool {
	// Simple directory copy implementation
	entries, err := os.ReadDir(src)
	if err != nil {
		return false
	}
	
	os.MkdirAll(dst, 0755)
	
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		
		if entry.IsDir() {
			if !copyDir(srcPath, dstPath) {
				return false
			}
		} else {
			if !copyFile(srcPath, dstPath) {
				return false
			}
		}
	}
	
	return true
}

func copyFile(src, dst string) bool {
	srcFile, err := os.Open(src)
	if err != nil {
		return false
	}
	defer srcFile.Close()
	
	dstFile, err := os.Create(dst)
	if err != nil {
		return false
	}
	defer dstFile.Close()
	
	_, err = io.Copy(dstFile, srcFile)
	return err == nil
}
