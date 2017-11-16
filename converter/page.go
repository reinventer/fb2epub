package converter

import (
	"encoding/xml"

	"fmt"
	"github.com/reinventer/fb2epub/epub"
	"strings"
)

func (c *FB2Converter) addPage(decoder *fb2decoder, e *epub.Epub, bodyName string, links map[string]string) error {
	var (
		err            error
		content        string
		currentSection string
		sectionDepth   int
		sectionsNum    int

		updatePage = func() error {
			if len(content) > 0 {
				c.chapter++
				currentSection = fmt.Sprintf("%s%d", bodyName, c.chapter)
				if err = e.AddPage(currentSection, "text", content); err != nil {
					return err
				}
				content = ""
			}
			sectionDepth = 0
			sectionsNum = 0

			return nil
		}
	)

	updatePage()

CYCLE:
	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		switch t := token.(type) {
		case xml.StartElement:
			var id string
			for _, a := range t.Attr {
				if a.Name.Local == "id" {
					id = ` id="` + a.Value + `"`
					break
				}
			}

			switch t.Name.Local {
			case "section":
				if sectionDepth == 0 {
					sectionsNum++
					if sectionsNum >= c.sectionsPerPage {
						if err = updatePage(); err != nil {
							return err
						}
					}
				}
				sectionDepth++

				fallthrough

			case "epigraph", "poem", "stanza", "text-author", "title", "subtitle", "cite":
				content += `<div class="` + t.Name.Local + `"` + id + `>`

			case "emphasis":
				content += "<em" + id + ">"

			case "p":
				content += "<p" + id + ">"

			case "v":
				content += `<p class="v"` + id + `>`

			case "empty-line":
				content += `<br />`

			case "a":
				var link string
				for _, a := range t.Attr {
					if a.Name.Local == "href" {
						if page, ok := links[a.Value]; ok {
							link = page + ".xhtml" + a.Value

						} else {
							link = a.Value
						}
						break
					}
				}
				content += `<a href="` + link + `"` + id + ">"

			case "image":
				for _, a := range t.Attr {
					if a.Name.Local == "href" {
						image := strings.TrimLeft(a.Value, "#")
						if len(image) > 0 {
							content += `<div class="image"><img src="` + image + `"` + id + " /></div>"
						}
						break
					}
				}
			}

		case xml.CharData:
			content += string(t)

		case xml.EndElement:
			switch t.Name.Local {
			case "section":
				sectionDepth--
				fallthrough

			case "epigraph", "poem", "stanza", "text-author", "title", "subtitle":
				content += `</div>`
			case "emphasis":
				content += `</em>`
			case "p", "v":
				content += `</p>`
			case "a":
				content += `</a>`
			case "body":
				break CYCLE

			}
		}
	}

	return updatePage()
}
