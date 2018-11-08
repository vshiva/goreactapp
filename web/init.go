package web

//go:generate ./gen.sh

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-contrib/static"
)

type binaryFS struct {
	fs http.FileSystem
}

func (b *binaryFS) Open(name string) (f http.File, err error) {
	f, err = b.fs.Open(name)
	if err == nil {
		return
	}
	f, err = b.fs.Open("index.html")
	return
}

func (b *binaryFS) Exists(prefix string, filepath string) bool {
	// always return true to server all SPA routes
	return true
}

// FS is the Web site file system. This serves the complete console binary
func FS() static.ServeFileSystem {
	fs := &assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "dist"}
	return &binaryFS{
		fs,
	}
}
