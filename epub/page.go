package epub

// AddPage adds page to epub
func (e *Epub) AddPage(pageName, guideType, content string) error {
	e.manifest += `<item id="` + pageName + `" href="` + pageName + `.xhtml" media-type="application/xhtml+xml" />` + "\n"
	e.spine += `<itemref idref="` + pageName + `" />` + "\n"

	data := struct {
		Title   string
		Content string
		Type    string
	}{
		Title:   e.title,
		Content: content,
		Type:    guideType,
	}

	return e.parseTemplate("page.tmpl", "OEBPS/"+pageName+".xhtml", data)
}
