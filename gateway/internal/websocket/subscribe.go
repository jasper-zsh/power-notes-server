package websocket

import (
	"github.com/mitchellh/mapstructure"
	"github.com/zeromicro/go-zero/core/logx"
	"powernotes-server/gateway/internal/model"
	"powernotes-server/gateway/internal/util"
)

var FileBroadcaster, ProjectBroadcaster *Broadcaster

func init() {
	FileBroadcaster = NewBroadcaster()
	ProjectBroadcaster = NewBroadcaster()
}

type OpenFileRequest struct {
	ProjectName string `mapstructure:"project_name"`
	FilePath    string `mapstructure:"file_path"`
}

func OnOpenFile(conn *WSConn, body map[string]interface{}) error {
	req := OpenFileRequest{}
	err := mapstructure.Decode(body, &req)
	if err != nil {
		return err
	}
	FileBroadcaster.Subscribe(conn, util.FileKey(req.ProjectName, req.FilePath))

	notes := make([]model.Note, 0)
	result := model.DB.Where("project_name = ? AND file_path = ?", req.ProjectName, req.FilePath).Find(&notes)
	if result.Error != nil {
		return result.Error
	}
	for _, note := range notes {
		err = PushToClient(conn, "note", note)
		if err != nil {
			logx.Errorf("WARN: Failed to push note to %s during open project", conn.ID)
			continue
		}
	}

	return nil
}

type CloseFileRequest struct {
	ProjectName string `mapstructure:"project_name"`
	FilePath    string `mapstructure:"file_path"`
}

func OnCloseFile(conn *WSConn, body map[string]interface{}) error {
	req := CloseFileRequest{}
	err := mapstructure.Decode(body, &req)
	if err != nil {
		return err
	}
	FileBroadcaster.Unsubscribe(conn, util.FileKey(req.ProjectName, req.FilePath))
	return nil
}

type ProjectRequest struct {
	ProjectName string `mapstructure:"project_name"`
}

func OnOpenProject(conn *WSConn, body map[string]interface{}) error {
	req := ProjectRequest{}
	err := mapstructure.Decode(body, &req)
	if err != nil {
		return err
	}
	ProjectBroadcaster.Subscribe(conn, req.ProjectName)
	flows := make([]model.Flow, 0)
	result := model.DB.Where("project_name = ?", req.ProjectName).Find(&flows)
	if result.Error != nil {
		return result.Error
	}
	for _, flow := range flows {
		err = PushToClient(conn, "flow", flow)
		if err != nil {
			logx.Errorf("WARN: Failed to push flow to %s during open project", conn.ID)
			continue
		}
		rels := make([]model.FlowNoteRelation, 0)
		result = model.DB.Where("flow_id = ?", flow.ID).Find(&rels)
		if result.Error != nil {
			logx.Errorf("Failed to get rels for flow %d", flow.ID)
			continue
		}
		for _, rel := range rels {
			err = PushToClient(conn, "flow_note_relation", rel)
			if err != nil {
				logx.Errorf("Failed to push flow note relation to %s", conn.ID)
				continue
			}
		}
	}
	notes := make([]*model.Note, 0)
	result = model.DB.Where("project_name = ?", req.ProjectName).Find(&notes)
	if result.Error != nil {
		return result.Error
	}
	for _, note := range notes {
		PushToClient(conn, "note", note)
	}
	return nil
}

func OnCloseProject(conn *WSConn, body map[string]interface{}) error {
	req := ProjectRequest{}
	err := mapstructure.Decode(body, &req)
	if err != nil {
		return err
	}
	ProjectBroadcaster.Unsubscribe(conn, req.ProjectName)
	return nil
}
