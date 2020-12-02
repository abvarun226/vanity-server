package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Rules for mapping go-import.
type Rules struct {
	ImportRoot string `json:"importroot"`
	VCS        string `json:"vcs"`
	RepoRoot   string `json:"reporoot"`
}

// ReloadRules is the handler that reloads mapping rules from config.
func (h *Handler) ReloadRules(w http.ResponseWriter, r *http.Request) {
	if err := h.GetImportRules(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("ok"))
}

// GetImportRules loads config file (import path mapping rules)
func (h *Handler) GetImportRules() error {
	data, err := ioutil.ReadFile(h.opts.ConfigFile)
	if err != nil {
		return ErrReadConfigFile
	}

	var rules []Rules
	if err = json.Unmarshal(data, &rules); err != nil {
		return ErrParseConfigFile
	}

	h.Lock()
	h.opts.MappingRules = rules
	h.Unlock()
	return nil
}

// custom errors.
var (
	ErrReadConfigFile  = fmt.Errorf("failed to read rules config file")
	ErrParseConfigFile = fmt.Errorf("failed to parse rules config file")
)
