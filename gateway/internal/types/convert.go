package types

import (
	"powernotes-server/gateway/internal/model"
)

func (api *Note) ToModel() *model.Note {
	return &model.Note{
		ID:            api.ID,
		ProjectName:   api.ProjectName,
		FileName:      api.FileName,
		FilePath:      api.FilePath,
		LineNumber:    api.LineNumber,
		EndLineNumber: api.EndLineNumber,
		Text:          api.Text,
		CreatedAt:     api.CreatedAt,
		UpdatedAt:     api.UpdatedAt,
	}
}

func NewNoteFromModel(m *model.Note) *Note {
	return &Note{
		ID:            m.ID,
		ProjectName:   m.ProjectName,
		FileName:      m.FileName,
		FilePath:      m.FilePath,
		LineNumber:    m.LineNumber,
		EndLineNumber: m.EndLineNumber,
		Text:          m.Text,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func (api *Flow) ToModel() *model.Flow {
	return &model.Flow{
		ID:          api.ID,
		ProjectName: api.Name,
		Name:        api.Name,
		CreatedAt:   api.CreatedAt,
		UpdatedAt:   api.UpdatedAt,
	}
}

func NewFlowFromModel(m *model.Flow) *Flow {
	return &Flow{
		ID:          m.ID,
		ProjectName: m.ProjectName,
		Name:        m.Name,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (api *FlowNoteRelation) ToModel() *model.FlowNoteRelation {
	return &model.FlowNoteRelation{
		FlowID:   api.FlowID,
		NoteID:   api.NoteID,
		Position: api.Position,
	}
}

func NewFlowNoteRelationFromModel(m *model.FlowNoteRelation) *FlowNoteRelation {
	return &FlowNoteRelation{
		FlowID:   m.FlowID,
		NoteID:   m.NoteID,
		Position: m.Position,
	}
}
