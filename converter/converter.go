package converter

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/reinventer/fb2epub/epub"
)

// FB2Converter - main object for converting fb2 to epub
type FB2Converter struct {
	fb2ReaderFunc   func() (io.ReadCloser, error)
	translit        bool
	sectionsPerPage int

	chapter int
}

// New returns converter
func New(fb2path string, sectionsPerPage int) (*FB2Converter, error) {
	var f func() (io.ReadCloser, error)

	if fb2path == "-" {
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}

		f = func() (io.ReadCloser, error) {
			return bytesReadCloser{
				Reader: bytes.NewReader(data),
			}, nil
		}

	} else {
		f = func() (io.ReadCloser, error) {
			return os.Open(fb2path)
		}

	}

	return &FB2Converter{
		fb2ReaderFunc:   f,
		sectionsPerPage: sectionsPerPage,
	}, nil
}

func (c *FB2Converter) links() (map[string]string, error) {
	decoder, err := c.decoder()
	if err != nil {
		return nil, err
	}

	defer decoder.close()

	var (
		links          = map[string]string{}
		currentBody    string
		currentSection string
		chapter        int
		sectionDepth   int
		sectionsNum    int

		updatePage = func() {
			sectionDepth = 0
			sectionsNum = 0
			chapter++
			currentSection = fmt.Sprintf("%s%d", currentBody, chapter)
		}
	)

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "section":
				if sectionDepth == 0 {
					sectionsNum++
					if sectionsNum >= c.sectionsPerPage {
						updatePage()
					}
				}
				sectionDepth++

			case "body":
				currentBody = "ch"
				for _, a := range t.Attr {
					if a.Name.Local == "name" && len(a.Value) > 0 {
						currentBody = a.Value
					}
				}
				updatePage()
			}

			for _, a := range t.Attr {
				if a.Name.Local == "id" && len(a.Value) > 0 {
					links[`#`+a.Value] = currentSection
					break
				}
			}
			////////

		case xml.EndElement:
			if t.Name.Local == "section" {
				sectionDepth--
			}

		}
	}

	return links, nil
}

func (c *FB2Converter) makeEpub(path string, links map[string]string, translit bool) error {
	decoder, err := c.decoder()
	if err != nil {
		return err
	}
	defer decoder.close()

	outEpub, err := epub.New(path, translit)
	if err != nil {
		return err
	}

	defer outEpub.Close()

	c.chapter = 0

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if t, ok := token.(xml.StartElement); ok {
			switch t.Name.Local {
			case "description":
				if err = c.fillDescription(decoder, outEpub); err != nil {
					return err
				}

			case "body":
				bodyName := "ch"
				for _, a := range t.Attr {
					if a.Name.Local == "name" && len(a.Value) > 0 {
						bodyName = a.Value
					}
				}

				if err = c.addPage(decoder, outEpub, bodyName, links); err != nil {
					return err
				}

			case "binary":
				content, err := decoder.getText()
				if err != nil {
					return err
				}

				var contentType, id string
				for _, a := range t.Attr {
					switch a.Name.Local {
					case "content-type":
						contentType = a.Value
					case "id":
						id = a.Value
					}
				}
				if err = outEpub.AddBinary(id, contentType, content); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Convert fb2 to epub
func (c *FB2Converter) Convert(epubPath string, translit bool) error {
	links, err := c.links()
	if err != nil {
		return err
	}

	return c.makeEpub(epubPath, links, translit)
}

type bytesReadCloser struct {
	*bytes.Reader
}

// Close - method to satisfy ReadCloser interface
func (b bytesReadCloser) Close() error {
	return nil
}
