package helper

import (
	"io"
	"net/http"
	"time"
)

// serveContent is a helper function that serves content using http.ServeContent.
// It handles range requests, sets the correct MIME type, and manages cache headers.
func serveContent(w http.ResponseWriter, req *http.Request, name string, modTime time.Time, content io.ReadSeeker) {
	// Set the appropriate MIME type (e.g., "audio/mpeg" for MP3 files)
	// You should set the MIME type based on the actual content you're serving
	mimeType := "audio/mpeg"
	w.Header().Set("Content-Type", mimeType)

	// Serve the content
	http.ServeContent(w, req, name, modTime, content)
}
