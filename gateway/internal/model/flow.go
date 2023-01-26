package model

type Flow struct {
	ID          int64  `gorm:"primaryKey,autoIncrement" json:"id" mapstructure:"id"`
	ProjectName string `json:"project_name" mapstructure:"project_name"`
	Name        string `json:"name" mapstructure:"name"`
	CreatedAt   int64  `json:"created_at" mapstructure:"created_at" gorm:"autoCreateTime:milli"`
	UpdatedAt   int64  `json:"updated_at" mapstructure:"updated_at" gorm:"autoUpdateTime:milli"`
}
