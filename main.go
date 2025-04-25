package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/makinori/jitsi-welcome/anime"
	"github.com/makinori/jitsi-welcome/common"
	"github.com/makinori/jitsi-welcome/jitsi"

	"github.com/charmbracelet/log"
)

var (
	//go:embed template.html assets
	staticContent embed.FS
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templateHTMLBytes, err := staticContent.ReadFile("template.html")
	if err != nil {
		log.Error("failed to get template", "err", err)
		http.Error(w, "failed to get template", http.StatusInternalServerError)
		return
	}

	if common.ConfigInDev {
		templateHTMLBytes, _ = os.ReadFile("template.html")
	}

	w.Write(templateHTMLBytes)

	// t, err := template.New("index.html").Parse(string(templateHTMLBytes))
	// if err != nil {
	// 	log.Error("failed to parse template", "err", err)
	// 	http.Error(w, "failed to parse template", http.StatusInternalServerError)
	// 	return
	// }

	// err = t.Execute(w, map[string]string{
	// 	"RandomRoomName": getJitsiRoomName(),
	// })

	// if err != nil {
	// 	log.Error("failed to run template", "err", err)
	// 	http.Error(w, "failed to run template", http.StatusInternalServerError)
	// 	return
	// }
}

func apiHandler(w http.ResponseWriter, r *http.Request, fn func() string) {
	data, err := json.Marshal(map[string]string{
		"name": fn(),
	})

	if err != nil {
		log.Error("failed to marshal json", "err", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func main() {
	assetsFS, err := fs.Sub(staticContent, "assets")
	if err != nil {
		log.Fatal("failed to get assets dir", "err", err)
	}

	http.Handle("GET /welcome/",
		http.StripPrefix("/welcome/", http.FileServerFS(assetsFS)),
	)

	http.HandleFunc("GET /welcome/anime-name", func(w http.ResponseWriter, r *http.Request) {
		apiHandler(w, r, anime.GenerateJitsiRoomName)
	})

	http.HandleFunc("GET /welcome/regular-name", func(w http.ResponseWriter, r *http.Request) {
		apiHandler(w, r, jitsi.GenerateRoomName)
	})

	http.HandleFunc("GET /{$}", indexHandler)

	addr := fmt.Sprintf(":%s", common.ConfigHTTPPort)

	log.Infof("listening at http://127.0.0.1%s", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("failed to start http server", "err", err)
	}
}
