package frontend

import (
	"embed"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// //go:embed dist/pikvm-automator/browser
var embeddedFiles embed.FS

func AddFrontend(mux *runtime.ServeMux) error {
	subFS, err := fs.Sub(embeddedFiles, "dist/pikvm-automator/browser")
	if err != nil {
		return fmt.Errorf("failed to create sub filesystem: %w", err)
	}

	system := http.FS(subFS)
	fileServer := http.FileServer(system)

	err = mux.HandlePath(http.MethodGet, "/*", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		if r.URL.Path != "/" {
			_, err = subFS.Open(strings.TrimPrefix(path.Clean(r.URL.Path), "/"))
			if err != nil {
				r.URL.Path = "/"
			}
		}
		fileServer.ServeHTTP(w, r)
	})
	if err != nil {
		return fmt.Errorf("failed to handle path: %w", err)
	}

	return nil
}
