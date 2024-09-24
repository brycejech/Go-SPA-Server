package internal

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime"
	"path/filepath"
	"sync"
)

func MustCreateFileCache(embeddedContent embed.FS, envReplacements map[string]string) FileCache {
	cache := &sync.Map{}

	contentFS, err := fs.Sub(embeddedContent, "artifact")
	if err != nil {
		log.Fatal(err)
	}

	fs.WalkDir(contentFS, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(fmt.Errorf("error walking embedded fs: %w", err))
		}

		if d.IsDir() {
			return nil
		}

		var (
			file      fs.File
			fileBytes []byte
		)

		if file, err = contentFS.Open(p); err != nil {
			log.Fatal(fmt.Errorf("error opening embdded file '%s': %w", p, err))
		}
		defer file.Close()

		if fileBytes, err = io.ReadAll(file); err != nil {
			log.Fatal(fmt.Errorf("error reading embedded file '%s': %w", p, err))
		}

		for old, new := range envReplacements {
			fileBytes = bytes.Replace(fileBytes, []byte(old), []byte(new), -1)
		}

		p = filepath.ToSlash(p)
		cache.Store(p, newCachedFile(p, fileBytes, getContentType(p)))

		return nil
	})

	return &fileCache{
		cache: cache,
	}
}

type fileCache struct {
	cache *sync.Map
}

func (c *fileCache) GetFile(filepath string) (CachedFile, error) {
	f, ok := c.cache.Load(filepath)
	if !ok {
		return nil, fmt.Errorf("file '%s' not found", filepath)
	}

	cachedFile, ok := f.(*cachedFile)
	if !ok {
		return nil, fmt.Errorf("error casting '%s' to cachedFile", f)
	}

	return cachedFile, nil
}

func getContentType(filename string) (contentType string) {
	ext := filepath.Ext(filename)
	if contentType = mime.TypeByExtension(ext); contentType == "" {
		contentType = "application/octect-stream"
	}
	return contentType
}
