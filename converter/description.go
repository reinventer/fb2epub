package converter

import (
	"encoding/xml"

	"github.com/reinventer/fb2epub/epub"
	"strings"
)

func (c *FB2Converter) fillDescription(decoder *fb2decoder, e *epub.Epub) error {
	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "title-info":
				if err = c.fillTitleInfo(decoder, e); err != nil {
					return err
				}
			}

		case xml.EndElement:
			if t.Name.Local == "description" {
				return nil
			}
		}
	}

	return nil
}

func (c *FB2Converter) fillTitleInfo(decoder *fb2decoder, e *epub.Epub) error {
	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "genre":
				genre, err := decoder.getText()
				if err != nil {
					return err
				}
				e.AddMetadataSubject(genre)

			case "author":
				firstName, middleName, lastname, err := c.getNames(decoder, "author")
				if err != nil {
					return err
				}
				e.AddMetadataAuthor(firstName, middleName, lastname)

			case "book-title":
				title, err := decoder.getText()
				if err != nil {
					return err
				}
				e.AddMetadataTitle(title)

			case "lang":
				lang, err := decoder.getText()
				if err != nil {
					return err
				}
				e.AddMetadataLanguage(lang)

			case "coverpage":
				if err = c.addCoverPage(decoder, e); err != nil {
					return err
				}
			}

		case xml.EndElement:
			if t.Name.Local == "title-info" {
				return nil
			}
		}
	}

	return nil
}

func (c *FB2Converter) getNames(decoder *fb2decoder, tag string) (fn, mn, ln string, err error) {
	var token xml.Token

	for {
		if token, err = decoder.Token(); err != nil {
			return
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "first-name":
				if fn, err = decoder.getText(); err != nil {
					return
				}

			case "middle-name":
				if mn, err = decoder.getText(); err != nil {
					return
				}

			case "last-name":
				if ln, err = decoder.getText(); err != nil {
					return
				}
			}

		case xml.EndElement:
			if t.Name.Local == tag {
				return
			}
		}
	}

	return
}

func (c *FB2Converter) addCoverPage(decoder *fb2decoder, e *epub.Epub) error {
	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == "image" {
				for _, a := range t.Attr {
					if a.Name.Local == "href" {
						imageName := a.Value
						if imageName[0] == '#' {
							imageName = strings.TrimLeft(imageName, "#")
						}
						if err = e.AddPage("cover", "cover", `<div class="cover"><img class="coverimage" alt="Cover" src="`+imageName+`" /></div>`); err != nil {
							return err
						}
						e.AddMetadataCover(imageName)
					}
				}
			}

		case xml.EndElement:
			if t.Name.Local == "coverpage" {
				return nil
			}
		}
	}
}
