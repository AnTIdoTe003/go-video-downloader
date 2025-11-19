# YouTube Downloader API Server

A REST API server built with **Gin** framework and Go that uses the `gostreampuller` package to download YouTube videos and provide metadata.

## Features

- 📋 **Get Video Metadata**: Fetch comprehensive video information without downloading
- 🔗 **Direct Download URLs**: Get direct YouTube download URLs in metadata response for use with other tools
- 📥 **Download Videos**: Download YouTube videos in various formats and resolutions
- 🌐 **CORS Enabled**: Ready for frontend integration
- ⚡ **Built with Gin**: Fast and efficient HTTP framework

## API Endpoints

### GET `/api/metadata?url=<youtube_url>`
Get metadata for a YouTube video without downloading it.

**Example:**
```bash
curl "http://localhost:8080/api/metadata?url=https://www.youtube.com/watch?v=dQw4w9WgXcQ"
```

**Response:**
```json
{
  "success": true,
  "metadata": {
    "id": "...",
    "title": "...",
    "duration": 212,
    "view_count": 1234567,
    ...
  },
  "download_url": "https://rr5---sn-xxx.googlevideo.com/videoplayback?..."
}
```

**Note**: The `download_url` field contains a direct download URL from YouTube that you can use with tools like `wget`, `curl`, or any other downloader. This URL is temporary and expires after some time.

### POST `/api/download`
Download a YouTube video directly to your local machine. This endpoint streams the file directly to your browser, triggering an automatic download.

**Request Body:**
```json
{
  "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
  "format": "mp4",
  "resolution": "720",
  "codec": "avc1"
}
```

**Response:**
- The video file is streamed directly as a binary download
- Your browser will automatically save it to your Downloads folder
- The filename is based on the video title

**Example with curl:**
```bash
curl -X POST "http://localhost:8080/api/download" \
  -H "Content-Type: application/json" \
  -d '{"url":"https://www.youtube.com/watch?v=dQw4w9WgXcQ","format":"mp4","resolution":"720"}' \
  --output video.mp4
```

### POST `/api/download-info`
Get download information and metadata without actually downloading the video.

**Request Body:**
```json
{
  "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
  "format": "mp4",
  "resolution": "720",
  "codec": "avc1"
}
```

**Response:**
```json
{
  "success": true,
  "file_path": "Video_Title.mp4",
  "metadata": { ... }
}
```

### GET `/health`
Health check endpoint.

## Installation

1. Make sure you have Go 1.24+ installed
2. Install dependencies:
```bash
cd youtube-api-server
go mod download
```

3. Run the server:
```bash
go run main.go
```

The server will start on port 8080 (or the port specified in the `PORT` environment variable).

## Configuration

- **Port**: Set `PORT` environment variable (default: 8080)
- **Temp Directory**: Videos are temporarily saved to `./temp_downloads/` during download, then automatically deleted after streaming

## Notes

- The first run will automatically install `yt-dlp` and `ffmpeg` if not already installed
- Videos are streamed directly to your browser - they download to your local machine's Downloads folder
- Temporary files are automatically cleaned up after streaming
- The API includes CORS headers for frontend integration
- Filenames are based on the video title (sanitized for filesystem compatibility)

