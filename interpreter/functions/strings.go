package functions

import (
	"strings"
	"unicode"
)

// toCamelCase converts a string to camelCase.
func ToCamelCase(s string) string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return ""
	}

	// Convert the first word to lowercase
	words[0] = strings.ToLower(words[0])

	// Capitalize the first letter of the remaining words
	for i := 1; i < len(words); i++ {
		words[i] = Capitalize(words[i])
	}

	return strings.Join(words, "")
}

// capitalize capitalizes the first letter of a word.
func Capitalize(word string) string {
	if len(word) == 0 {
		return word
	}
	runes := []rune(word)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// lowerFirst converts the first letter of the first word in the input string to lowercase.
func LowerFirst(s string) string {
	if len(s) == 0 {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// lowerWords converts all words in the input string to lowercase.
func LowerWords(s string) string {
	return strings.ToLower(s)
}

// upFirst capitalizes the first letter of the first word in the input string.
func UpFirst(s string) string {
	if len(s) == 0 {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func UpWords(s string) string {
	return strings.Title(s)
}

// ToSnakeCase converts a string to snake_case.
func ToSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			// Add an underscore before uppercase letters (except at the start)
			if i > 0 {
				result = append(result, '_')
			}
			// Convert the letter to lowercase
			r = unicode.ToLower(r)
		} else if r == ' ' || r == '-' || r == '.' || r == ',' {
			// Replace spaces and certain punctuation with underscores
			r = '_'
		}
		result = append(result, r)
	}
	return string(result)
}

// ToKebabCase converts a string to kebab-case.
func ToKebabCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			// Add a hyphen before uppercase letters (except at the start)
			if i > 0 {
				result = append(result, '-')
			}
			// Convert the letter to lowercase
			r = unicode.ToLower(r)
		} else if r == ' ' || r == '_' || r == '.' || r == ',' {
			// Replace spaces and certain punctuation with hyphens
			r = '-'
		}
		result = append(result, r)
	}
	return string(result)
}

// ToPascalCase converts a string to PascalCase.
func ToPascalCase(s string) string {
	var result []rune
	capitalizeNext := true

	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if capitalizeNext {
				r = unicode.ToUpper(r)
				capitalizeNext = false
			} else {
				r = unicode.ToLower(r)
			}
			result = append(result, r)
		} else {
			capitalizeNext = true
		}
	}
	return string(result)
}

// ToDotCase converts a string to dot.case.
func ToDotCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if i > 0 && (unicode.IsUpper(r) || (!unicode.IsLetter([]rune(s)[i-1]) && unicode.IsLetter(r))) {
				result = append(result, '.')
			}
			result = append(result, unicode.ToLower(r))
		} else if len(result) > 0 && result[len(result)-1] != '.' {
			result = append(result, '.')
		}
	}
	if len(result) > 0 && result[len(result)-1] == '.' {
		result = result[:len(result)-1] // Remove trailing dot, if any
	}
	return string(result)
}
