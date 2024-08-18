package frontend

import (
	"embed"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

//go:embed dist/pikvm-automator/browser
var embeddedFiles embed.FS

func AddFrontend(mux *runtime.ServeMux) {
	subFS, err := fs.Sub(embeddedFiles, "dist/pikvm-automator/browser")
	if err != nil {
		panic(err)
	}

	system := http.FS(subFS)
	fileServer := http.FileServer(system)

	mux.HandlePath(http.MethodGet, "/*", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		if r.URL.Path != "/" {
			_, err = subFS.Open(strings.TrimPrefix(path.Clean(r.URL.Path), "/"))
			if err != nil {
				r.URL.Path = "/"
			}
		}
		fileServer.ServeHTTP(w, r)
	})
}
