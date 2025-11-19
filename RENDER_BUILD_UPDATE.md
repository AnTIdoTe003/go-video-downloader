# Render Build Command with yt-dlp Update

Since pre-deploy commands are only available for paid Render plans, we need to update yt-dlp in the build command itself.

## Updated Build Command

Use this build command in Render:

```bash
go mod download && go build -o main && (yt-dlp -U || ~/.gostreampuller/bin/yt-dlp -U || true)
```

This will:
1. Download Go dependencies
2. Build the application
3. Try to update yt-dlp (if available in PATH or local install)
4. Continue even if update fails (the `|| true` ensures build doesn't fail)

## Alternative: Simpler Build Command

If the above doesn't work, use this simpler version:

```bash
go mod download && go build -o main
```

The application will automatically update yt-dlp in the background on first run (non-blocking).

## How It Works

The code now includes:
1. **Auto-update on startup**: When the app starts and yt-dlp is already installed, it will try to update it in the background
2. **Non-blocking**: The update happens in a goroutine, so it won't slow down your first request
3. **Silent failure**: If the update fails, the app continues with the existing version

## Manual Update (If Needed)

If you need to manually update yt-dlp on Render:

1. SSH into your Render instance (if available)
2. Run: `~/.gostreampuller/bin/yt-dlp -U`
3. Or restart the service to trigger auto-update

## Current Render Settings

**Build Command:**
```bash
go mod download && go build -o main
```

**Start Command:**
```bash
./main
```

The auto-update will happen in the background when the app starts.

