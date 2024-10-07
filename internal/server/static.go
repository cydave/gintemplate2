package server

import (
	"io/fs"
	"log"

	"github.com/cydave/gintemplate2/internal/assets"
)

// Get root assets in the static FS. Every file that is in the top-level
// directory is returned.
func getRootAssets() []string {
	ret := make([]string, 0)
	entries, err := fs.ReadDir(assets.Static, "static")
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ret = append(ret, "/"+entry.Name())
	}
	return ret
}
