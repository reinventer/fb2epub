package epub

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"text/template"

	"github.com/google/uuid"
	"github.com/reinventer/fb2epub/fs"
)

// Epub structure to represent epub file
type Epub struct {
	out      *zip.Writer
	uuid     string
	title    string
	language string
	metadata string
	manifest string
	spine    string

	templates *template.Template
	translit  bool
}

// New returns new Epub object
func New(outFile string, translit bool) (*Epub, error) {
	var (
		w      io.WriteCloser
		err    error
		result = &Epub{
			uuid: uuid.New().String(),
			manifest: `<item id="css" href="main.css" media-type="text/css" />` + "\n" +
				`<item href="fonts/Lora-Bold.ttf" media-type="application/octet-stream" id="font1"/>` + "\n" +
				`<item href="fonts/Lora-BoldItalic.ttf" media-type="application/octet-stream" id="font2"/>` + "\n" +
				`<item href="fonts/Lora-Italic.ttf" media-type="application/octet-stream" id="font3"/>` + "\n" +
				`<item href="fonts/Lora-Regular.ttf" media-type="application/octet-stream" id="font4"/>` + "\n",
			translit: translit,
		}
	)

	if outFile == "-" {
		w = os.Stdout
	} else if w, err = os.Create(outFile); err != nil {
		return nil, err
	}

	result.out = zip.NewWriter(w)

	tmpls, err := fs.AssetDir("tmpl")
	if err != nil {
		return nil, err
	}

	result.templates = template.New("")
	for _, tmpl := range tmpls {
		content, err := fs.Asset("tmpl/" + tmpl)
		if err != nil {
			return nil, err
		}

		result.templates, err = result.templates.New(tmpl).Parse(string(content))
		if err != nil {
			return nil, err
		}
	}

	for _, route := range [][]string{
		{"files/mimetype", "mimetype"},
		{"files/container.xml", "META-INF/container.xml"},
		{"files/main.css", "OEBPS/main.css"},
		{"files/LiberationSerif-Bold.ttf", "OEBPS/fonts/LiberationSerif-Bold.ttf"},
		{"files/LiberationSerif-BoldItalic.ttf", "OEBPS/fonts/LiberationSerif-BoldItalic.ttf"},
		{"files/LiberationSerif-Italic.ttf", "OEBPS/fonts/LiberationSerif-Italic.ttf"},
		{"files/LiberationSerif-Regular.ttf", "OEBPS/fonts/LiberationSerif-Regular.ttf"},
	} {
		if err = result.copyFile(route[0], route[1]); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// Close zip writer
func (e *Epub) Close() error {
	var err error

	if e.out == nil {
		return errors.New("writer not initialized")
	}

	if err = e.addContent(); err != nil {
		return err
	}

	return e.out.Close()
}

func (e *Epub) addContent() error {
	data := struct {
		UUID     string
		Title    string
		Language string
		Metadata string
		Manifest string
		Spine    string
	}{
		UUID:     e.uuid,
		Title:    e.title,
		Language: e.language,
		Metadata: e.metadata,
		Manifest: e.manifest,
		Spine:    e.spine,
	}

	if err := e.parseTemplate("content.tmpl", "OEBPS/content.opf", data); err != nil {
		return err
	}
	return e.parseTemplate("toc.tmpl", "OEBPS/toc.ncx", data)
}

func (e *Epub) parseTemplate(tmplName, outFile string, data interface{}) error {
	f, err := e.out.Create(outFile)
	if err != nil {
		return err
	}

	return e.templates.ExecuteTemplate(f, tmplName, data)
}

func (e *Epub) copyFile(source, destination string) error {
	out, err := e.out.Create(destination)
	if err != nil {
		return err
	}

	content, err := fs.Asset(source)
	if err != nil {
		return err
	}

	_, err = out.Write(content)
	return err
}
