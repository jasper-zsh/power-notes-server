package service

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"powernotes-server/handler"
	"powernotes-server/model"
)

const (
	EventRemoveFlowNoteRelation = "remove_flow_note_relation"
	EventFlowNoteRelation       = "flow_note_relation"
)

func SaveFlowNoteRelation(rel *model.FlowNoteRelation) (*model.FlowNoteRelation, error) {
	flow := model.Flow{}
	result := model.DB.Where("id = ?", rel.FlowID).First(&flow)
	if result.Error != nil {
		return nil, result.Error
	}
	notify := make([]*model.FlowNoteRelation, 0)
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		origins := make([]*model.FlowNoteRelation, 0)
		result = tx.Where("flow_id = ? AND position = ?", rel.FlowID, rel.Position).Find(&origins)
		if result.Error != nil {
			return result.Error
		}
		if len(origins) > 0 {
			// 重排序
			result = tx.Model(&model.FlowNoteRelation{}).Where("flow_id = ? AND position >= ?", flow.ID, rel.Position).Update("position", gorm.Expr("position + 1"))
			if result.Error != nil {
				return result.Error
			}
			result = tx.Where("flow_id = ? AND position > ?", flow.ID, rel.Position).Find(&origins)
			if result.Error != nil {
				return result.Error
			}
			notify = append(notify, origins...)
		}
		result = tx.Clauses(clause.OnConflict{UpdateAll: true}).Save(rel)
		if result.Error != nil {
			return result.Error
		}
		notify = append(notify, rel)
		return nil
	})
	if err != nil {
		return nil, err
	}
	for _, r := range notify {
		handler.ProjectBroadcaster.Broadcast(flow.ProjectName, EventFlowNoteRelation, r)
	}
	return rel, nil
}

func RemoveFlowNoteRelation(flowID int64, noteID int64) error {
	flow := model.Flow{}
	result := model.DB.Where("id = ?", flowID).First(&flow)
	if result.Error != nil {
		return result.Error
	}
	result = model.DB.Where("flow_id = ? AND note_id = ?", flowID, noteID).Delete(&model.FlowNoteRelation{})
	if result.Error != nil {
		return result.Error
	}
	handler.ProjectBroadcaster.Broadcast(flow.ProjectName, EventRemoveFlowNoteRelation, &model.FlowNoteRelation{
		FlowID: flowID,
		NoteID: noteID,
	})
	return nil
}

type SwapFlowNote struct {
	FlowID   int64 `json:"flow_id"`
	NoteID   int64 `json:"note_id"`
	Position int   `json:"position"`
	Offset   int   `json:"offset"`
}

func SwapFlowNoteRelation(req *SwapFlowNote) error {
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		flow := model.Flow{}
		result := tx.Where("id = ?", req.FlowID).First(&flow)
		if result.Error != nil {
			return result.Error
		}
		src := model.FlowNoteRelation{}
		result = tx.Where("flow_id = ? AND note_id = ? AND position = ?", req.FlowID, req.NoteID, req.Position).First(&src)
		if result.Error != nil {
			return result.Error
		}
		positionOp := "<"
		order := "position "
		offset := req.Offset
		if req.Offset > 0 {
			positionOp = ">"
			order += "ASC"
		} else {
			order += "DESC"
			offset = -req.Offset
		}
		origins := make([]*model.FlowNoteRelation, 0)
		result = tx.Where(fmt.Sprintf("flow_id = ? AND position %s ?", positionOp), req.FlowID, req.Position).Order(order).Limit(offset).Find(&origins)
		if result.Error != nil {
			return nil
		}
		if len(origins) > 0 {
			target := origins[0]
			result = tx.Model(&model.FlowNoteRelation{}).Where("flow_id = ? AND note_id = ?", src.FlowID, src.NoteID).Update("position", target.Position)
			if result.Error != nil {
				return result.Error
			}
			result = tx.Model(&model.FlowNoteRelation{}).Where("flow_id = ? AND note_id = ?", target.FlowID, target.NoteID).Update("position", src.Position)
			if result.Error != nil {
				return result.Error
			}
			notify := []model.FlowNoteRelation{
				{FlowID: src.FlowID, NoteID: src.NoteID, Position: target.Position},
				{FlowID: target.FlowID, NoteID: target.NoteID, Position: src.Position},
			}
			for _, rel := range notify {
				handler.ProjectBroadcaster.Broadcast(flow.ProjectName, "flow_note_relation", rel)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
