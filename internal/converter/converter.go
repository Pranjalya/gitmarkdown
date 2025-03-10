package converter

type Converter interface {
	Supports(filePath string) bool
	Convert(filePath string) (string, error)
	GetLanguage(filePath string) string
}

func GetConverter(filePath string, converters []Converter, defaultConverter Converter) Converter {
	for _, conv := range converters {
		if conv.Supports(filePath) {
			return conv
		}
	}
	return defaultConverter
}
