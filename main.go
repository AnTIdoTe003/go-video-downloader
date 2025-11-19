package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"youtube-api-server/pkg/downloader"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type MetadataResponse struct {
	Success        bool                      `json:"success"`
	Metadata       *downloader.VideoMetadata `json:"metadata,omitempty"`
	DownloadURL    string                    `json:"download_url,omitempty"` // Direct YouTube download URL
	Error          string                    `json:"error,omitempty"`
}

type DownloadResponse struct {
	Success      bool                      `json:"success"`
	DownloadURL  string                    `json:"download_url,omitempty"`
	FilePath     string                    `json:"file_path,omitempty"` // Expected filename
	Metadata     *downloader.VideoMetadata `json:"metadata,omitempty"`
	Error        string                    `json:"error,omitempty"`
}

type DownloadRequest struct {
	URL        string `json:"url"`
	Format     string `json:"format,omitempty"`     // mp4, webm, etc.
	Resolution string `json:"resolution,omitempty"` // 720, 1080, etc.
	Codec      string `json:"codec,omitempty"`      // avc1, vp9, etc.
}

// Store for temporary downloaded files (cleaned up after streaming)
var tempDir = "./temp_downloads"

func init() {
	// Create temp directory if it doesn't exist
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		log.Printf("Warning: Could not create temp directory: %v", err)
	}
}

func main() {
	// Set Gin to release mode (optional, for production)
	// gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:3001"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// API routes
	api := router.Group("/api")
	{
		api.GET("/metadata", getMetadataHandler)
		api.POST("/download", downloadStreamHandler)
		api.POST("/download-info", downloadInfoHandler)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 YouTube Downloader API Server running on port %s\n", port)
	fmt.Printf("📁 Temp directory: %s\n", tempDir)
	fmt.Printf("💡 Downloads will stream directly to your browser\n")
	log.Fatal(router.Run(":" + port))
}

func getMetadataHandler(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(400, MetadataResponse{
			Success: false,
			Error:   "URL parameter is required",
		})
		return
	}

	// Validate YouTube URL
	if !isValidYouTubeURL(url) {
		c.JSON(400, MetadataResponse{
			Success: false,
			Error:   "Invalid YouTube URL",
		})
		return
	}

	// Fetch metadata
	metadata, err := downloader.GetVideoMetadata(url)
	if err != nil {
		c.JSON(500, MetadataResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to fetch metadata: %v", err),
		})
		return
	}

	// Get direct download URL from YouTube
	downloadURL, err := getDirectDownloadURL(url)
	if err != nil {
		log.Printf("Warning: Could not get direct download URL: %v", err)
		// Continue without download URL - metadata is still useful
	}

	c.JSON(200, MetadataResponse{
		Success:     true,
		Metadata:    metadata,
		DownloadURL: downloadURL,
	})
}

// getDirectDownloadURL gets the direct download URL from YouTube using yt-dlp
func getDirectDownloadURL(url string) (string, error) {
	// Ensure binaries are installed
	if err := ensureBinariesInstalled(); err != nil {
		return "", fmt.Errorf("failed to ensure binaries are installed: %w", err)
	}

	// Get yt-dlp path - we need to find it
	ytdlpPath := findYTDLPPath()
	if ytdlpPath == "" {
		return "", fmt.Errorf("yt-dlp not found")
	}

	// Use yt-dlp with -g flag to get direct URL
	// -g: Print video URL instead of downloading
	// -f best: Get best quality format
	cmd := exec.Command(ytdlpPath, "-g", "-f", "best", url)

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("yt-dlp error: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("failed to execute yt-dlp: %w", err)
	}

	// yt-dlp may return multiple URLs (video + audio), take the first one
	urls := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(urls) > 0 && urls[0] != "" {
		return urls[0], nil
	}

	return "", fmt.Errorf("no download URL found")
}

// findYTDLPPath finds the yt-dlp binary path
func findYTDLPPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "yt-dlp" // Fallback to PATH
	}

	binDir := filepath.Join(homeDir, ".gostreampuller", "bin")
	candidates := []string{
		filepath.Join(binDir, "yt-dlp"),
		filepath.Join(binDir, "yt-dlp.exe"),
		"yt-dlp", // System PATH
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
		if path == "yt-dlp" {
			// Check if it's in PATH
			if _, err := exec.LookPath("yt-dlp"); err == nil {
				return "yt-dlp"
			}
		}
	}

	return "yt-dlp" // Final fallback
}

// ensureBinariesInstalled ensures yt-dlp and ffmpeg are installed
func ensureBinariesInstalled() error {
	// This is a simplified version - the actual implementation is in the downloader package
	// We'll just check if yt-dlp exists
	ytdlpPath := findYTDLPPath()
	if ytdlpPath == "" || ytdlpPath == "yt-dlp" {
		// Try to find it in PATH
		if _, err := exec.LookPath("yt-dlp"); err != nil {
			return fmt.Errorf("yt-dlp not found. Please install it or let the package auto-install on first download")
		}
	}
	return nil
}

// downloadStreamHandler streams the video directly to the client, triggering browser download
func downloadStreamHandler(c *gin.Context) {
	var req DownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err)})
		return
	}

	if req.URL == "" {
		c.JSON(400, gin.H{"error": "URL is required"})
		return
	}

	// Validate YouTube URL
	if !isValidYouTubeURL(req.URL) {
		c.JSON(400, gin.H{"error": "Invalid YouTube URL"})
		return
	}

	// Set defaults
	if req.Format == "" {
		req.Format = "mp4"
	}
	if req.Resolution == "" {
		req.Resolution = "720"
	}
	if req.Codec == "" {
		req.Codec = "avc1"
	}

	// Fetch metadata first to get video title for filename
	metadata, err := downloader.GetVideoMetadata(req.URL)
	var filename string
	if err == nil && metadata != nil {
		// Use video title as filename (sanitized)
		filename = sanitizeFilename(metadata.Title) + "." + req.Format
	} else {
		// Fallback to timestamp-based filename
		filename = fmt.Sprintf("video_%d.%s", time.Now().UnixNano(), req.Format)
	}

	// Download video to temp directory
	filePath, err := downloader.DownloadVideoToDir(
		req.URL,
		req.Format,
		req.Resolution,
		req.Codec,
		tempDir,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to download video: %v", err)})
		return
	}

	// Clean up temp file after streaming
	defer func() {
		if err := os.Remove(filePath); err != nil {
			log.Printf("Warning: Failed to clean up temp file %s: %v", filePath, err)
		}
	}()

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to open file: %v", err)})
		return
	}
	defer file.Close()

	// Get file info for Content-Length
	fileInfo, err := file.Stat()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get file info: %v", err)})
		return
	}

	// Set headers to trigger browser download
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	c.Header("Content-Transfer-Encoding", "binary")

	// Stream the file to the client
	c.DataFromReader(200, fileInfo.Size(), "application/octet-stream", file, nil)
}

// downloadInfoHandler returns metadata and download info without actually downloading
func downloadInfoHandler(c *gin.Context) {
	var req DownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, DownloadResponse{
			Success: false,
			Error:   fmt.Sprintf("Invalid request body: %v", err),
		})
		return
	}

	if req.URL == "" {
		c.JSON(400, DownloadResponse{
			Success: false,
			Error:   "URL is required",
		})
		return
	}

	// Validate YouTube URL
	if !isValidYouTubeURL(req.URL) {
		c.JSON(400, DownloadResponse{
			Success: false,
			Error:   "Invalid YouTube URL",
		})
		return
	}

	// Set defaults
	if req.Format == "" {
		req.Format = "mp4"
	}
	if req.Resolution == "" {
		req.Resolution = "720"
	}
	if req.Codec == "" {
		req.Codec = "avc1"
	}

	// Fetch metadata
	metadata, err := downloader.GetVideoMetadata(req.URL)
	if err != nil {
		c.JSON(500, DownloadResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to fetch metadata: %v", err),
		})
		return
	}

	// Generate expected filename
	var filename string
	if metadata != nil {
		filename = sanitizeFilename(metadata.Title) + "." + req.Format
	} else {
		filename = fmt.Sprintf("video_%d.%s", time.Now().UnixNano(), req.Format)
	}

	c.JSON(200, DownloadResponse{
		Success:  true,
		Metadata: metadata,
		FilePath: filename, // Expected filename
	})
}

// sanitizeFilename removes invalid characters from filename
func sanitizeFilename(name string) string {
	// Remove invalid characters for filenames
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", "\n", "\r"}
	result := name
	for _, char := range invalidChars {
		result = strings.ReplaceAll(result, char, "_")
	}
	// Limit length
	if len(result) > 100 {
		result = result[:100]
	}
	return strings.TrimSpace(result)
}

func isValidYouTubeURL(url string) bool {
	validPatterns := []string{
		"youtube.com/watch",
		"youtu.be/",
		"youtube.com/embed/",
		"youtube.com/v/",
		"youtube.com/shorts/",
	}

	urlLower := strings.ToLower(url)
	for _, pattern := range validPatterns {
		if strings.Contains(urlLower, pattern) {
			return true
		}
	}
	return false
}
