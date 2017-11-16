package epub

import "encoding/base64"

// AddBinary adds binary file to epub
func (e *Epub) AddBinary(id, contentType, base64Content string) error {
	e.manifest += `<item id="` + id + `" href="` + id + `" media-type="` + contentType + `" />` + "\n"
	data, err := base64.StdEncoding.DecodeString(base64Content)
	if err != nil {
		return err
	}

	f, err := e.out.Create("OEBPS/" + id)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	return err
}
