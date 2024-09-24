package internal

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type FileCache interface {
	GetFile(filename string) (CachedFile, error)
}

type CachedFile interface {
	Path() string
	Data() []byte
	ContentType() string
	ContentLength() string
}

func NewSPAServer(port string, cache FileCache, staticDir string) *http.Server {
	handler := &spaHandler{
		cache:     cache,
		staticDir: staticDir,
	}

	server := &http.Server{
		Handler:      http.HandlerFunc(handler.handlerFunc),
		Addr:         fmt.Sprintf("localhost:%v", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return server
}

type spaHandler struct {
	cache     FileCache
	staticDir string
}

func (s *spaHandler) handlerFunc(w http.ResponseWriter, r *http.Request) {
	reqFile := strings.TrimPrefix(r.URL.Path, "/")

	var (
		f    CachedFile
		err  error
		err2 error
	)

	if f, err = s.cache.GetFile(reqFile); err != nil {
		hasStaticDir := len(s.staticDir) > 0
		if hasStaticDir && strings.HasPrefix(reqFile, s.staticDir) {
			// Caught error in static dir, should 404
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
			return
		}
		if f, err2 = s.cache.GetFile("index.html"); err2 != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
			return
		}
	}

	w.Header().Add("content-type", f.ContentType())
	w.Header().Add("content-length", f.ContentLength())
	w.WriteHeader(http.StatusOK)
	w.Write(f.Data())
}
