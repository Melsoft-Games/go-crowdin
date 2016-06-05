package crowdin

type File struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	NodeType        string `json:"node_type"`
	Phrases         string `json:"phrases"`
	Translated      string `json:"translated"`
	Approved        string `json:"approved"`
	Words           string `json:"words"`
	WordsTranslated string `json:"words_translated"`
	WordsApproved   string `json:"words_approved"`
}

type Files struct {
	Files []File `json:"files"`
}