package handlers

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/labstack/echo/v4"
)

type DocsHandler struct{}

func NewDocsHandler() *DocsHandler {
	return &DocsHandler{}
}

var rootDir string

func init() {
	rootDir, _ = os.Getwd()
}

func createDocs(name string, c echo.Context) error {
	filename := path.Join(rootDir, "docs", "schema", fmt.Sprintf("%s.html", name))
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		c.Logger().Errorf("failed to parse template: %v", err)
		return err
	}
	var bufferHTML bytes.Buffer
	if err := tmpl.Execute(&bufferHTML, nil); err != nil {
		c.Logger().Errorf("failed to execute template: %v", err)
		return err
	}
	return c.HTML(200, bufferHTML.String())
}

func (h *DocsHandler) Docs(c echo.Context) error {
	switch c.Param("path") {
	case "openapi.yaml":
		return c.File(path.Join(rootDir, "docs", "schema", "openapi.yaml"))
	case "redoc":
		return createDocs("redoc", c)
	case "swagger":
		return createDocs("swagger", c)
	case "scalar":
		return createDocs("scalar", c)
	default:
		return echo.NewHTTPError(404, fmt.Sprintf("Document %s not found", c.Param("name")))
	}
}
