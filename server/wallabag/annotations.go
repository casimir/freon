package wallabag

type Annotation struct {
	ID                     int               `json:"id"`
	AnnotatorSchemaVersion int               `json:"annotator_schema_version"`
	CreatedAt              *Time             `json:"created_at"`
	UpdatedAt              *Time             `json:"updated_at"`
	Quote                  *string           `json:"quote"`
	Ranges                 []AnnotationRange `json:"ranges"`
	Text                   string            `json:"text"`
	User                   *string           `json:"user"`
}

type AnnotationRange struct {
	Start       *string   `json:"start"`
	End         *string   `json:"end"`
	StartOffset *MagicInt `json:"startOffset"`
	EndOffset   *MagicInt `json:"endOffset"`
}
