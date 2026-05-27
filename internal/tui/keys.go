package tui

import tea "github.com/charmbracelet/bubbletea"

// qwertyToRussian maps QWERTY runes to their Russian (ЙЦУКЕН) layout counterparts.
var qwertyToRussian = map[rune]rune{
	'q': 'й', 'w': 'ц', 'e': 'у', 'r': 'к', 't': 'е', 'y': 'н', 'u': 'г', 'i': 'ш', 'o': 'щ', 'p': 'з',
	'[': 'х', ']': 'ъ', 'a': 'ф', 's': 'ы', 'd': 'в', 'f': 'а', 'g': 'п', 'h': 'р', 'j': 'о', 'k': 'л',
	'l': 'д', ';': 'ж', '\'': 'э', 'z': 'я', 'x': 'ч', 'c': 'с', 'v': 'м', 'b': 'и', 'n': 'т', 'm': 'ь',
	',': 'б', '.': 'ю', '/': '.', '`': 'ё',
	'Q': 'Й', 'W': 'Ц', 'E': 'У', 'R': 'К', 'T': 'Е', 'Y': 'Н', 'U': 'Г', 'I': 'Ш', 'O': 'Щ', 'P': 'З',
	'{': 'Х', '}': 'Ъ', 'A': 'Ф', 'S': 'Ы', 'D': 'В', 'F': 'А', 'G': 'П', 'H': 'Р', 'J': 'О', 'K': 'Л',
	'L': 'Д', ':': 'Ж', '"': 'Э', 'Z': 'Я', 'X': 'Ч', 'C': 'С', 'V': 'М', 'B': 'И', 'N': 'Т', 'M': 'Ь',
	'<': 'Б', '>': 'Ю', '?': ',', '~': 'Ё',
}

// keyMatches reports whether msg corresponds to the expected rune, taking into
// account common alternative keyboard layouts (currently Russian ЙЦУКЕN).
func keyMatches(msg tea.KeyMsg, expected rune) bool {
	if len(msg.Runes) != 1 {
		return false
	}
	r := msg.Runes[0]
	if r == expected {
		return true
	}
	if ru, ok := qwertyToRussian[expected]; ok && r == ru {
		return true
	}
	return false
}
