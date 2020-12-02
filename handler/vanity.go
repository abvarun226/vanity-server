package handler

import (
	"bytes"
	"net/http"
	"strings"
	"text/template"
)

// VanityServer redirects browsers to godoc or go tool to VCS repository.
func (h *Handler) VanityServer(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		status := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(status), status)
		return
	}

	pkgName := r.Host + r.URL.Path

	// If go-get param is absent, redirect to godoc URL.
	if r.FormValue("go-get") != "1" {
		if h.opts.GodocURL == "" {
			w.Write([]byte(nothingHere))
			return
		}
		url := h.opts.GodocURL + pkgName
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}

	// go-import mapping rules.
	var importRoot, vcs, repoRoot string
	for _, rule := range h.opts.MappingRules {
		if strings.HasPrefix(strings.ToLower(pkgName), strings.ToLower(rule.ImportRoot)+"/") {
			repoName := strings.Replace(strings.ToLower(pkgName), strings.ToLower(rule.ImportRoot), "", -1)
			repoName = strings.Split(repoName, "/")[1]

			importRoot = rule.ImportRoot + "/" + repoName
			repoRoot = rule.RepoRoot + "/" + repoName
			vcs = rule.VCS

			break
		}
	}

	// Create HTML template with go-import <meta> tag
	d := struct {
		ImportRoot string
		VCS        string
		RepoRoot   string
	}{
		ImportRoot: importRoot,
		VCS:        vcs,
		RepoRoot:   repoRoot,
	}

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, &d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=900")
	w.Write(buf.Bytes())
}

var tmpl = template.Must(template.New("main").Parse(`<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="go-import" content="{{.ImportRoot}} {{.VCS}} {{.RepoRoot}}">
</head>
</html>
`))

const nothingHere = `
<html>
<head></head>
<body><h5>Nothing here. Move along.</h5></body>
</html>
`
