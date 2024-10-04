package server

import (
	"embed"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type staticFS struct {
	root http.FileSystem
}

func (s *staticFS) Open(name string) (http.File, error) {
	f, err := s.root.Open(name)
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return nil, os.ErrNotExist
	}
	return f, nil
}

func ServeStaticFS(prefix string, static embed.FS) gin.HandlerFunc {
	sfs := &staticFS{http.FS(static)}
	return func(c *gin.Context) {
		file := c.Param("filepath")
		fp := filepath.Join(prefix, filepath.Clean(file))
		f, err := sfs.Open(fp)
		if err != nil {
			c.Writer.WriteHeader(http.StatusNotFound)
			return
		}
		f.Close()
		http.FileServer(sfs).ServeHTTP(c.Writer, c.Request)
	}
}
