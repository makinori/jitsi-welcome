package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"regexp"
	"text/template"

	"github.com/makinori/jitsi-welcome/anime"
	"github.com/makinori/jitsi-welcome/common"

	"github.com/charmbracelet/log"
)

var (
	//go:embed template.html assets
	staticContent embed.FS
)

func getJitsiRoomName() string {
	name, err := anime.GetRandomAnimeName(common.ConfigAniListUsername)
	if err != nil {
		log.Error("failed to get random anime name", "err", err)
		return "FailedToGetRandomRoomName"
	}

	// accents and diacritics included
	r := regexp.MustCompile("(?i)[^a-zA-Z\u00C0-\u024F\u1E00-\u1EFF]")

	// jitsi supports special characters though
	// name = norm.NFKD.String(name)

	name = r.ReplaceAllString(name, "")

	return name
}

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

	t, err := template.New("index.html").Parse(string(templateHTMLBytes))
	if err != nil {
		log.Error("failed to parse template", "err", err)
		http.Error(w, "failed to parse template", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, map[string]string{
		"RandomRoomName": getJitsiRoomName(),
	})

	if err != nil {
		log.Error("failed to run template", "err", err)
		http.Error(w, "failed to run template", http.StatusInternalServerError)
		return
	}
}

func main() {
	assetsFS, err := fs.Sub(staticContent, "assets")
	if err != nil {
		log.Fatal("failed to get assets dir", "err", err)
	}

	http.Handle("GET /welcome-assets/",
		http.StripPrefix("/welcome-assets/", http.FileServerFS(assetsFS)),
	)

	http.HandleFunc("GET /{$}", indexHandler)

	addr := fmt.Sprintf(":%s", common.ConfigHTTPPort)

	log.Infof("listening at http://127.0.0.1%s", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("failed to start http server", "err", err)
	}
}
