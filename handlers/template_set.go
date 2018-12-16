package handlers

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
)

const (
	templateExt = ".html"
)

// TemplateSet -
type TemplateSet struct {
	reload bool
	root   string
	*template.Template
}

// TemplateSetConfig -
type TemplateSetConfig struct {
	Root         string
	AlwaysReload bool
}

// NewTemplateSet -
func NewTemplateSet(cfg *TemplateSetConfig) (*TemplateSet, error) {
	return makeTemplateSet(cfg.Root, cfg.AlwaysReload)
}

// MustTemplateSet -
func MustTemplateSet(cfg *TemplateSetConfig) *TemplateSet {
	ts, err := NewTemplateSet(cfg)
	if err != nil {
		panic(err)
	}
	return ts
}

// ExecuteTemplate wrap inner template.Template.ExecuteTemplate() for provider reload utility.
func (ts *TemplateSet) ExecuteTemplate(w io.Writer, name string, data interface{}) error {
	if ts.reload {
		return ts.executeLatestTemplate(w, name, data)
	}
	return ts.Template.ExecuteTemplate(w, name, data)
}

func (ts *TemplateSet) executeLatestTemplate(w io.Writer, name string, data interface{}) error {
	// TODO need lock
	refreshed, err := findTemplateFiles(ts.root)
	if err != nil {
		return err
	}

	t, err := newTemplate(refreshed)
	if err != nil {
		return fmt.Errorf("failed to parse template files: %s", err)
	}

	return t.ExecuteTemplate(w, name, data)
}

//cSpell:words tmpls
func makeTemplateSet(root string, alwaysParse bool) (*TemplateSet, error) {
	tmpls, err := findTemplateFiles(root)
	if err != nil {
		return nil, err
	}
	t, err := newTemplate(tmpls)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template files: %s", err)
	}

	return &TemplateSet{Template: t, root: root, reload: alwaysParse}, nil
}

func findTemplateFiles(root string) ([]string, error) {
	var tmpls []string
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if isTemplateFile(path) {
			tmpls = append(tmpls, path)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to walk template dirs: %s", err)
	}
	return tmpls, nil
}

func isTemplateFile(path string) bool {
	return strings.HasSuffix(path, templateExt)
}

func newTemplate(files []string) (*template.Template, error) {
	return template.New("templateSet").Funcs(funcMap).ParseFiles(files...)
}

var funcMap = template.FuncMap{
	"formatTime":        tfFormatTime,
	"formatPostMetaKey": tfFormatPostMetaID,
	"formatBytes":       tfFormatBytes,
	"iterate":           tfIterate,
}

func tfFormatTime(t time.Time) string {
	return t.Format("2006-01-02:15:04:05")
}

func tfFormatPostMetaID(k *datastore.Key) string {
	return strconv.FormatInt(k.ID, 10)
}

func tfFormatBytes(b []byte) string {
	return string(b)
}

func tfIterate(n int) []int {
	var ns = make([]int, 0, n)
	for i := 0; i < n; i++ {
		ns[i] = i
	}
	return ns
}
