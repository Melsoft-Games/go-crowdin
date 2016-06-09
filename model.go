package crowdin

// Used for AddFile() API call
type AddFileOptions struct {
	// Note: Used only when uploading CSV (or XLS/XLSX) file to define data columns mapping.
	// Acceptable value is the combination of the following constants:
	// "identifier" — Column contains string identifier.
	// "source_phrase" — Column contains only source string (in result string will contain same string).
	// "source_or_translation" — Column contains source string but when exporting same column should contain translation (also when uploading existing translations, the value from this column will be used as a translated string).
	// "translation" — Column contains translated string (when imported file already contains translations).
	// "context" — Column contains some comments on source string. Context information.
	// "max_length" — Column contains max. length of translation for this string.
	// "none" — Do not import column.
	Scheme string

	// Used when uploading CSV (or XLS/XLSX) files via API. Defines whether first line should be imported or it contains columns headers. May not contain value.
	FirstLineContainsHeader bool

	// Files array that should be added to Crowdin project. Array keys should contain file names with path in Crowdin project.
	Files map[string]string

	// Acceptable values are:
	// empty value or "auto" — Try to detect file type by extension or MIME type
	// "gettext" — GNU GetText (*.po, *.pot)
	// "qtts" — Nokia Qt (*.ts)
	// "dklang" — Delphi DKLang (*.dklang)
	// "android" — Android (*.xml)
	// "resx" — .NET (*.resx)
	// "properties" — Java (*.properties)
	// "macosx" — Mac OS X / iOS (*.strings)
	// "blackberry" — BlackBerry (*.rrc)
	// "Symbian" — Symbian (*.lXX)
	// "flex" — Adobe Flex (*.properties)
	// "bada" — Samsung Bada (*.xml)
	// "txt" — Plain Text (*.txt)
	// "srt" — SubRip .srt (*.srt)
	// "sbv" — Youtube .sbv (*.sbv)
	// "xliff" — XLIFF (*.xliff)
	// "html" — HTML (*.html, *.htm, *.xhtml, *.xhtm)
	// "dtd" — Mozilla DTD (*.dtd)
	// "chrome" — Google Chrome Extension (*.json)
	// "yaml" — Ruby On Rails (*.yaml)
	// "csv" — Comma Separated Values (*.csv)
	// "rc" — Windows Resources (*.rc)
	// "wxl" — WiX Installer Resources (*.wxl)
	// "nsh" — NSIS Installer Resources (*.nsh)
	// "joomla" — Joomla localizable resources (*.ini)
	// "ini" — Generic INI (*.ini)
	// "isl" — ISL (*.isl)
	// "resw" — Windows 8 Metro (*.resw)
	// "resjson" — Windows 8 Metro (*.resjson)
	// "docx" — Microsoft Office and OpenOffice.org Documents (*.docx, *.dotx, *.odt, *.ott, *.xslx, *.xltx, *.pptx, *.potx, *.ods, *.ots, *.odg, *.otg, *.odp, *.otp, *.idml)
	// "md" — Markdown (*.md, *.text, *.markdown...)
	// "mediawiki" — MediaWiki (*.wiki, *.wikitext, *.mediawiki)
	// "play" — Play Framework
	// "haml" — Haml (*.haml)
	// "arb" — Application Resource Bundle (*.arb)
	// "vtt" — Video Subtitling and WebVTT (*.vtt)
	Type string
}

type UpdateFileOptions struct {
	// Note: Used only when uploading CSV (or XLS/XLSX) file to define data columns mapping.
	// Acceptable value is the combination of the following constants:
	// "identifier" — Column contains string identifier.
	// "source_phrase" — Column contains only source string (in result string will contain same string).
	// "source_or_translation" — Column contains source string but when exporting same column should contain translation (also when uploading existing translations, the value from this column will be used as a translated string).
	// "translation" — Column contains translated string (when imported file already contains translations).
	// "context" — Column contains some comments on source string. Context information.
	// "max_length" — Column contains max. length of translation for this string.
	// "none" — Do not import column.
	Scheme string

	// Used when uploading CSV (or XLS/XLSX) files via API. Defines whether first line should be imported or it contains columns headers. May not contain value.
	FirstLineContainsHeader bool

	// Files array that should be added to Crowdin project. Array keys should contain file names with path in Crowdin project.
	Files map[string]string
}

type UploadTranslationsOptions struct {
	// Target language. With a single call it's possible to upload translations for several files but only into one of the languages.
	Language string

	// Translated files array. Array keys should contain file names in Crowdin.
	Files map[string]string

	// Defines whether to add translation if there is the same translation previously added. Acceptable values are: 0 or 1. Default is 0.
	ImportDuplicates string
}

type responseLanguageStatus struct {
	Files []struct {
		ID              string `json:"id"`
		Name            string `json:"name"`
		NodeType        string `json:"node_type"`
		Phrases         string `json:"phrases"`
		Translated      string `json:"translated"`
		Approved        string `json:"approved"`
		Words           string `json:"words"`
		WordsTranslated string `json:"words_translated"`
		WordsApproved   string `json:"words_approved"`
	} `json:"files"`
}

type responseAddFile struct {
	Success bool `json:"success"`
	Stats   struct {
		Files []struct {
			FileID  int    `json:"file_id"`
			Name    string `json:"name"`
			Strings int    `json:"strings"`
			Words   int    `json:"words"`
		} `json:"files"`
	} `json:"stats"`
}

type responseUploadTranslation struct {
	Success bool `json:"success"`
	Stats   struct {
		Files []struct {
			Name   string `json:"name"`
			Status string `json:"status"`
		} `json:"files"`
	} `json:"stats"`
}

type responseGeneral struct {
	Success bool `json:"success"`
}
