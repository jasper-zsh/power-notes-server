package model

type Note struct {
	ID            int64  `gorm:"primaryKey;autoIncrement" json:"id" mapstructure:"id"`
	ProjectName   string `json:"project_name" mapstructure:"project_name" gorm:"type:varchar;size:255;uniqueIndex:idx_note"`
	FileName      string `json:"file_name" mapstructure:"file_name"`
	FilePath      string `json:"file_path" mapstructure:"file_path" gorm:"type:varchar;size:255;uniqueIndex:idx_note"`
	LineNumber    int    `json:"line_number" mapstructure:"line_number" gorm:"uniqueIndex:idx_note"`
	EndLineNumber int    `json:"end_line_number" mapstructure:"end_line_number" gorm:"uniqueIndex:idx_note"`
	Text          string `json:"text" mapstructure:"text"`
	CreatedAt     int64  `json:"created_at" mapstructure:"created_at" gorm:"autoCreateTime:milli"`
	UpdatedAt     int64  `json:"updated_at" mapstructure:"updated_at" gorm:"autoUpdateTime:milli"`
}
