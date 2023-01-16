package service

import (
	"gorm.io/gorm"
	"powernotes-server/model"
)

func RemoveFlow(id int64) (*model.Flow, error) {
	flow := model.Flow{}
	result := model.DB.Where("id = ?", id).First(&flow)
	if result.Error != nil {
		return nil, result.Error
	}
	err := model.DB.Transaction(func(tx *gorm.DB) error {
		result = tx.Delete(&model.Flow{}, id)
		if result.Error != nil {
			return result.Error
		}
		result = tx.Where("flow_id = ?", id).Delete(&model.FlowNoteRelation{})
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &flow, nil
}
