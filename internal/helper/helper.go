package helper

import (
	"io"
	"net/http"
	"time"
)

func ServeContent(w http.ResponseWriter, req *http.Request, name string, modTime time.Time, content io.ReadSeeker) {
	mimeType := "audio/mpeg"
	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// Serve the content
	http.ServeContent(w, req, name, modTime, content)
}
