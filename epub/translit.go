package epub

import (
	"unicode/utf8"
)

func transliterate(s string) string {
	var (
		trMap = map[rune]string{
			'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo", 'ж': "zh", 'з': "z", 'и': "i", 'й': "j",
			'к': "k", 'л': "l", 'м': "m", 'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u", 'ф': "f",
			'х': "h", 'ц': "c", 'ч': "ch", 'ш': "sh", 'щ': "sch", 'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
			'А': "A", 'Б': "B", 'В': "V", 'Г': "G", 'Д': "D", 'Е': "E", 'Ё': "Yo", 'Ж': "Zh", 'З': "Z", 'И': "I", 'Й': "J",
			'К': "K", 'Л': "L", 'М': "M", 'Н': "N", 'О': "O", 'П': "P", 'Р': "R", 'С': "S", 'Т': "T", 'У': "U", 'Ф': "F",
			'Х': "H", 'Ц': "C", 'Ч': "Ch", 'Ш': "Sh", 'Щ': "Sch", 'Ъ': "", 'Ы': "Y", 'Ь': "", 'Э': "E", 'Ю': "Yu", 'Я': "Ya",
		}

		source = []byte(s)
		result = make([]byte, 0, len(s))
	)

	for len(source) > 0 {
		r, size := utf8.DecodeRune(source)

		if replace, ok := trMap[r]; ok {
			result = append(result, []byte(replace)...)
		} else {
			result = append(result, source[:size]...)
		}

		source = source[size:]
	}

	return string(result)
}
