package svc

import (
	"gorm.io/gorm/clause"
	"powernotes-server/gateway/internal/model"
	"powernotes-server/gateway/internal/websocket"
)

func SaveNote(note *model.Note) (*model.Note, error) {
	result := model.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "project_name"}, {Name: "file_path"}, {Name: "line_number"}, {Name: "end_line_number"}},
		UpdateAll: true,
	}).Save(&note)
	if result.Error != nil {
		return nil, result.Error
	}
	_ = websocket.ProjectBroadcaster.Broadcast(note.ProjectName, "note", &note)
	return note, nil
}

func RemoveNote(id int64) (*model.Note, error) {
	note := model.Note{}
	result := model.DB.Where("id = ?", id).First(&note)
	if result.Error != nil {
		return nil, result.Error
	}
	result = model.DB.Delete(&model.Note{}, id)
	if result.Error != nil {
		return nil, result.Error
	}
	rels := make([]*model.FlowNoteRelation, 0)
	result = model.DB.Where("note_id = ?", id).Find(&rels)
	if result.Error != nil {
		return nil, result.Error
	}
	result = model.DB.Where("note_id = ?", id).Delete(&model.FlowNoteRelation{})
	if result.Error != nil {
		return nil, result.Error
	}
	_ = websocket.ProjectBroadcaster.Broadcast(note.ProjectName, "remove_note", note)
	for _, rel := range rels {
		_ = websocket.ProjectBroadcaster.Broadcast(note.ProjectName, EventRemoveFlowNoteRelation, rel)
	}
	return &note, nil
}
