# YouTube Bot Detection Fix

## What Was Fixed

Added browser-like headers to all yt-dlp commands to bypass YouTube's bot detection:

- **User-Agent**: Chrome browser user agent
- **Referer**: YouTube domain
- **Accept-Language**: English language headers
- **Accept**: HTML/XML accept headers

These headers make yt-dlp requests look like they're coming from a real browser instead of a bot.

## Changes Made

Updated all yt-dlp command calls in:
1. `GetVideoMetadataWithContext()` - Metadata fetching
2. `DownloadVideoToDirWithProgress()` - Video downloads
3. `DownloadAudioToDirWithProgress()` - Audio downloads
4. `getDirectDownloadURL()` - Direct URL fetching

## If It Still Doesn't Work

### Option 1: Update yt-dlp

YouTube frequently changes their detection methods. Update yt-dlp:

```bash
# On Render, you can add this to a pre-deploy command:
yt-dlp -U
```

Or set it in Render's **Pre-Deploy Command**:
```bash
yt-dlp -U || true
```

### Option 2: Use Cookies (Advanced)

If headers aren't enough, you can use cookies:

1. Export cookies from your browser using a browser extension
2. Save cookies to a file (e.g., `cookies.txt`)
3. Add to yt-dlp commands: `--cookies cookies.txt`

### Option 3: Update yt-dlp Version

Make sure you're using the latest version. The auto-installer should get the latest, but you can force an update.

## Testing

After deploying, test with:
```bash
curl "https://your-app.onrender.com/api/metadata?url=https://www.youtube.com/watch?v=A6EobyJczEY"
```

## Monitoring

If you still get bot detection errors:
1. Check yt-dlp version: `yt-dlp --version`
2. Update yt-dlp: `yt-dlp -U`
3. Check yt-dlp GitHub for latest fixes: https://github.com/yt-dlp/yt-dlp



