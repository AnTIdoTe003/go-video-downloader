# Troubleshooting yt-dlp on Render

## Debug Endpoint

I've added a debug endpoint to check yt-dlp status. After deploying, visit:

```
https://your-app.onrender.com/debug/yt-dlp
```

This will show:
- yt-dlp path
- Home directory
- Binary directory
- Whether the binary exists
- Whether it's executable
- Version information

## Common Issues

### Issue 1: yt-dlp not found

**Symptoms**: Error like "yt-dlp: command not found" or "failed to execute yt-dlp"

**Solution**:
1. Check the debug endpoint to see where it's looking
2. The auto-installer should run on first API call
3. Check Render logs for installation messages

### Issue 2: Binary not executable

**Symptoms**: Permission denied errors

**Solution**: The code now automatically sets executable permissions after installation. If it still fails:
1. Check the debug endpoint
2. The binary should be at `~/.gostreampuller/bin/yt-dlp`

### Issue 3: Home directory issues

**Symptoms**: Can't create `.gostreampuller` directory

**Solution**: On Render, the home directory should be writable. If not:
1. Check Render logs for errors
2. The code falls back to system PATH if home directory fails

### Issue 4: Auto-installation not running

**Symptoms**: No installation messages in logs

**Possible causes**:
1. `GOSTREAMPULLER_NO_AUTO_INSTALL=1` is set (check environment variables)
2. Installation already attempted (check `installAttempted` flag)
3. Binaries found in system PATH

**Solution**:
1. Check environment variables in Render
2. Make a test API call to trigger installation
3. Check logs for any error messages

## Manual Installation (If Needed)

If auto-installation fails, you can manually install yt-dlp on Render:

1. SSH into your Render instance (if available)
2. Run:
   ```bash
   mkdir -p ~/.gostreampuller/bin
   curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o ~/.gostreampuller/bin/yt-dlp
   chmod +x ~/.gostreampuller/bin/yt-dlp
   ```

## Check Logs

After deploying, check Render logs for:
- `[gostreampuller] Installing yt-dlp...`
- `[gostreampuller] ✓ yt-dlp installed successfully`
- `[gostreampuller] Using yt-dlp at: ...`

If you see errors, they'll be in the logs.

## Test After Deployment

1. Visit: `https://your-app.onrender.com/debug/yt-dlp`
2. Check the response for yt-dlp status
3. Make a test API call: `https://your-app.onrender.com/api/metadata?url=https://www.youtube.com/watch?v=dQw4w9WgXcQ`
4. Check logs for installation messages



