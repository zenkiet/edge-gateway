package static

import (
	"edge-gateway/internal/pkg/response"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	immutableMaxAge  = time.Hour * 24 * 365 // 1 year
	revalidateMaxAge = time.Hour * 24       // 1 day
	lastModifiedFmt  = time.RFC1123
)

type Handler struct {
	manager *Manager
}

func NewHandler(m *Manager) *Handler {
	return &Handler{manager: m}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("pos_version")
	bundleValue := ""
	if err == nil {
		bundleValue = cookie.Value
	}

	version, ok := h.manager.GetVersion(bundleValue)
	if !ok {
		response.InternalServerError(w, "No version available")
		return
	}

	isLatest := h.manager.IsLatest(version.ID)
	reqPath := filepath.FromSlash(r.URL.Path)
	fullPath := filepath.Join(version.DirPath, reqPath)

	info, err := os.Stat(fullPath)
	if err == nil && !info.IsDir() {
		h.serveFile(w, r, fullPath, info, isLatest)
		return
	}

	// Fallback
	if strings.Contains(filepath.Base(reqPath), ".") && !strings.HasSuffix(reqPath, ".html") {
		response.InternalServerError(w, "Resource not found")
		return
	}

	indexHtml := filepath.Join(version.DirPath, "index.html")
	if info, err := os.Stat(indexHtml); err == nil {
		h.serveFile(w, r, indexHtml, info, isLatest)
		return
	}

	response.InternalServerError(w, "Index not found")
}

func (h *Handler) serveFile(w http.ResponseWriter, r *http.Request, path string, info os.FileInfo, isLatest bool) {
	etag := fmt.Sprintf(`"%s-%d"`, info.Name(), info.Size()+info.ModTime().Unix())
	if r.Header.Get("If-None-Match") == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	contentType := mime.TypeByExtension(filepath.Ext(path))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("ETag", etag)

	isIndexHTML := info.Name() == "index.html"

	switch {
	case isIndexHTML:
		w.Header().Set("Cache-Control", "no-cache, must-revalidate")
		w.Header().Set("Pragma", "no-cache")

	case isLatest:
		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", revalidateMaxAge))
		w.Header().Set("Last-Modified", info.ModTime().UTC().Format(lastModifiedFmt))

	default:
		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d, immutable", immutableMaxAge))
	}

	http.ServeFile(w, r, path)
}
