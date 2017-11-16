package epub

// AddMetadataSubject adds subject to metadata
func (e *Epub) AddMetadataSubject(s string) {
	e.metadata += "<dc:subject>" + e.transliterateIfNeed(s) + "</dc:subject>\n"
}

// AddMetadataAuthor adds author to metadata
func (e *Epub) AddMetadataAuthor(firstName, middleName, lastName string) {
	name := e.transliterateIfNeed(firstName)
	if len(middleName) != 0 {
		name += e.transliterateIfNeed(" " + middleName)
	}
	ln := e.transliterateIfNeed(lastName)

	e.metadata += `<dc:creator opf:file-as="` + ln + ", " + name + `" opf:role="aut" xmlns:opf="http://www.idpf.org/2007/opf">` + name + " " + ln + `</dc:creator>` + "\n"
}

// AddMetadataTitle adds title to metadata
func (e *Epub) AddMetadataTitle(s string) {
	e.title = s
	e.metadata += "<dc:title>" + e.transliterateIfNeed(s) + "</dc:title>\n"
}

// AddMetadataLanguage adds language to metadata
func (e *Epub) AddMetadataLanguage(s string) {
	e.language = s
	e.metadata += `<dc:language xsi:type="dcterms:RFC3066">` + s + `</dc:language>` + "\n"
}

// AddMetadataCover add cover to metadata
func (e *Epub) AddMetadataCover(imageName string) {
	e.metadata += `<meta name="cover" content="` + imageName + `" />` + "\n"
}

func (e *Epub) transliterateIfNeed(s string) string {
	if e.translit {
		return transliterate(s)
	}
	return s
}
