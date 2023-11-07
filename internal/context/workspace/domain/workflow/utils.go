package workflow

// SupportedLanguages get supported workflow languages
func SupportedLanguages() []string {
	return []string{
		string(LanguageWDL),
		string(LanguageNextflow),
	}
}
