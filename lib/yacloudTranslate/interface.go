package yacloudTranslate

// YaTranslate is the interface that wraps Yandex Cloud Translate API
//
// There are 3 methods available:
// DetectLanguage() - detects the language of the text
// ListLanguages() - retrieves the list of supported languages
// Translate() - translates the text to the specified language
type YaTranslate interface {
	DetectLanguage(req DetectLanguageRequest) (DetectLanguageResponse, error)
	ListLanguages(req ListLanguagesRequest) (ListLanguagesResponse, error)
	Translate(req TranslateRequest) (TranslateResponse, error)
}

// DetectLanguageRequest
type DetectLanguageRequest struct {
	// The text to detect the language for.
	// The maximum string length in characters is 1000.
	Text string `json:"text"`

	// List of the most likely languages. These languages will be given preference when detecting the text language. Specified in ISO 639-1 format (for example, ru).
	// To get the list of supported languages, use a listLanguages request.
	// The maximum number of elements is 10. The maximum string length in characters for each value is 3.
	LanguageCodeHints []string `json:"languageCodeHints"`
}

// DetectLanguageResponse
type DetectLanguageResponse struct {
	// The text language in ISO 639-1 format (for example, ru).
	// To get the language name, use a listLanguages request.
	LanguageCode string `json:"languageCode"`
}

type ListLanguagesRequest struct {
}

type Language struct {
	// The language code. Specified in ISO 639-1 format (for example, en).
	Code string `json:"code"`
	// The name of the language (for example, English).
	Name string `json:"name"`
}

type ListLanguagesResponse struct {
	// List of supported languages.
	Languages []Language `json:"languages"`
}

type Format string

const (
	// Text without markup. Default value.
	FormatPlainText Format = "PLAIN_TEXT"

	// Text in the HTML format.
	FormatHtml Format = "HTML"
)

type TranslateGlossaryConfig struct {
	// Pass glossary data in the request. Currently, only this way to pass glossary is supported.
	GlossaryData GlossaryData `json:"glossaryData"`
}

type GlossaryData struct {
	// Required. Array of text pairs.
	// The maximum total length of all source texts is 10000 characters. The maximum total length of all translated texts is 10000 characters.
	// The number of elements must be in the range 1-50.
	GlossaryPairs []GlossaryPair `json:"glossaryPairs"`
}

type GlossaryPair struct {
	// Required. Text in the source language.
	SourceText string `json:"sourceText"`

	// Required. Text in the target language.
	TranslatedText string `json:"translatedText"`

	Exact bool `json:"exact,omitempty"`
}

type TranslateRequest struct {
	// The text language to translate from. Specified in ISO 639-1 format (for example, ru).
	// Required for translating with glossary.
	// The maximum string length in characters is 3.
	SourceLanguageCode string `json:"sourceLanguageCode,omitempty"`

	// Required. The target language to translate the text. Specified in ISO 639-1 format (for example, en).
	// The maximum string length in characters is 3.
	TargetLanguageCode string `json:"targetLanguageCode"`

	// Format of the text.
	// FormatPlainText: Text without markup. Default value.
	// FormatHtml: Text in the HTML format.
	Format Format `json:"format,omitempty"`

	// Required. Array of the strings to translate. The maximum total length of all strings is 10000 characters.
	// Must contain at least one element.
	Texts []string `json:"texts"`

	// Do not specify this field, custom models are not supported yet.
	// The maximum string length in characters is 50.
	Model string `json:"model,omitempty"`

	// Glossary to be applied for the translation. For more information, see GlossaryConfig.
	GlossaryConfig *TranslateGlossaryConfig `json:"glossaryConfig,omitempty"`

	// Use speller
	Speller bool `json:"speller,omitempty"`
}

type TranslateResponse struct {
	// Array of the translations.
	Translations []TranslatedText `json:"translations"`
}

type TranslatedText struct {
	// Translated text.
	Text string `json:"text"`

	// The language code of the source text. Specified in ISO 639-1 format (for example, en).
	DetectedLanguageCode string `json:"detectedLanguageCode"`
}
