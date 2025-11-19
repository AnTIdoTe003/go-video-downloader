# Fix Render Build Command

## The Problem

Render is using the default build command: `go build -tags netgo -ldflags '-s -w' -o main`

This fails because:
1. It doesn't run `go mod download` first
2. The replace directive needs the parent directory to be available

## The Solution

### Step 1: Update Build Command in Render

Go to your Render service settings and change the **Build Command** to:

```bash
cd youtube-api-server && go mod download && go build -o main
```

Or if you want to keep the optimization flags:

```bash
cd youtube-api-server && go mod download && go build -tags netgo -ldflags '-s -w' -o main
```

### Step 2: Update Start Command

Make sure the **Start Command** is:

```bash
./main
```

### Step 3: Verify Root Directory

**CRITICAL**: Make sure **Root Directory** is set to `youtube-api-server` in Render settings.

If Root Directory is empty or set incorrectly, the `replace` directive won't work.

## Complete Render Configuration

**Settings:**
- **Root Directory**: `youtube-api-server`
- **Build Command**: `cd youtube-api-server && go mod download && go build -o main`
- **Start Command**: `./main`
- **Environment**: `Go`

**Environment Variables:**
- `PORT`: `8080` (optional)
- `GIN_MODE`: `release` (optional)

## Alternative: If Root Directory Can't Be Set

If you can't set Root Directory, use this build command instead:

```bash
go mod download && go build -o main
```

But you MUST deploy from the `youtube-api-server` directory, not the root.

## Verify go.sum is Committed

Before deploying, make sure `go.sum` is in your repository:

```bash
cd youtube-api-server
git add go.sum
git commit -m "Add go.sum"
git push
```

Then trigger a new deployment in Render.

