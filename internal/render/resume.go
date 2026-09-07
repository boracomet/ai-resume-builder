package render

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"strings"
	"sync"

	"github.com/boracomet/ai-resume-builder/internal/i18n"
	"github.com/boracomet/ai-resume-builder/internal/models"
	webassets "github.com/boracomet/ai-resume-builder/web"
)

type Renderer struct {
	tmpl *template.Template
	css  string
}

var (
	renderer     *Renderer
	rendererOnce sync.Once
	rendererErr  error
)

func Get() (*Renderer, error) {
	rendererOnce.Do(func() {
		renderer, rendererErr = newRenderer()
	})
	return renderer, rendererErr
}

func newRenderer() (*Renderer, error) {
	cssBytes, err := fs.ReadFile(webassets.Templates, "templates/resume.css")
	if err != nil {
		return nil, fmt.Errorf("read resume.css: %w", err)
	}

	funcMap := template.FuncMap{
		"fullURL":    fullURL,
		"displayURL": displayURL,
		"photoSrc":   photoSrc,
		"t":          labelForLang,
		"langCode":   i18n.NormalizeLang,
	}

	tmpl, err := template.New("resume.html").Funcs(funcMap).ParseFS(webassets.Templates, "templates/resume.html")
	if err != nil {
		return nil, fmt.Errorf("parse resume template: %w", err)
	}

	return &Renderer{
		tmpl: tmpl,
		css:  string(cssBytes),
	}, nil
}

func (r *Renderer) Render(profile *models.CVProfile) (string, error) {
	profile.Normalize()

	data := map[string]interface{}{
		"Profile": profile,
		"CSS":     template.CSS(r.css),
	}

	var buf bytes.Buffer
	if err := r.tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func fullURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return raw
	}
	if strings.HasPrefix(raw, "mailto:") {
		return raw
	}
	if strings.Contains(raw, "@") && !strings.Contains(raw, " ") {
		return "mailto:" + raw
	}
	return "https://" + strings.TrimPrefix(raw, "//")
}

func photoSrc(raw string) template.URL {
	return template.URL(strings.TrimSpace(raw))
}

func displayURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	raw = strings.TrimPrefix(raw, "https://")
	raw = strings.TrimPrefix(raw, "http://")
	raw = strings.TrimPrefix(raw, "mailto:")
	return raw
}

func labelForLang(lang, key string) string {
	return i18n.T(lang, key)
}
