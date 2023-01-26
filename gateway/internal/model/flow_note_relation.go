package model

type FlowNoteRelation struct {
	FlowID   int64 `gorm:"primaryKey" json:"flow_id" mapstructure:"flow_id"`
	NoteID   int64 `gorm:"primaryKey" json:"note_id" mapstructure:"note_id"`
	Position int   `json:"position" mapstructure:"position"`
}
