# Deployment Guide for Render

## Important: Deploy the Entire Repository

Since `youtube-api-server` uses a local `replace` directive pointing to the parent directory, you **must deploy the entire repository**, not just the `youtube-api-server` folder.

## Render Configuration Steps

### 1. Repository Setup
- Connect your **entire repository** (not just the `youtube-api-server` subdirectory)
- The repository should contain both `youtube-api-server/` and the parent `gostreampuller` package

### 2. Render Service Settings

**Basic Settings:**
- **Name**: `youtube-api-server` (or your preferred name)
- **Environment**: `Go`
- **Region**: Choose your preferred region
- **Branch**: `main` (or your default branch)

**Build & Deploy:**
- **Root Directory**: `youtube-api-server` ⚠️ **IMPORTANT: Set this!**
- **Build Command**: `go mod download && go build -o main`
- **Start Command**: `./main`

**Environment Variables:**
- `PORT`: `8080` (Render will override this automatically, but you can set a default)
- `GIN_MODE`: `release` (optional, for production mode)

### 3. Make Sure go.sum is Committed

Before deploying, ensure `go.sum` is committed to your repository:

```bash
cd youtube-api-server
go mod tidy
git add go.sum go.mod
git commit -m "Add go.sum for deployment"
git push
```

### 4. Deploy

Click "Create Web Service" or "Save Changes" in Render. The build should now work because:
- The entire repository is available
- The `replace` directive can find the parent directory
- `go.sum` is present for dependency verification

## Troubleshooting

### Error: "missing go.sum entry"

**Solution**: Make sure `go.sum` is committed to your repository:
```bash
cd youtube-api-server
go mod tidy
git add go.sum
git commit -m "Add go.sum"
git push
```

### Error: "cannot find module providing package"

**Solution**:
1. Verify **Root Directory** is set to `youtube-api-server` in Render settings
2. Make sure the entire repository is connected (not just a subdirectory)
3. The `replace` directive needs access to the parent directory

### Error: "unknown revision v0.0.0"

**Solution**: This happens if the replace directive can't find the parent directory. Make sure:
- Root Directory is set correctly
- The entire repository structure is present

### Build Command Alternative

If the default build command doesn't work, try:
```bash
cd youtube-api-server && go mod download && go mod verify && go build -o main
```

## Alternative: Copy Package Locally

If you can't deploy the entire repository, you can copy the required packages:

```bash
# In youtube-api-server directory
mkdir -p vendor
cp -r ../downloader vendor/
cp -r ../internal vendor/
```

Then update imports in `main.go` to use `./vendor/...` paths. However, this is not recommended as it duplicates code.

## Verify Deployment

After deployment, test the API:
```bash
curl https://your-service.onrender.com/health
```

Expected response:
```json
{"status":"ok"}
```
