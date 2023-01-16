package handler

import (
	"github.com/kataras/neffos"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"powernotes-server/model"
	"powernotes-server/util"
	"powernotes-server/websocket"
)

var FileBroadcaster, ProjectBroadcaster *websocket.Broadcaster

func init() {
	FileBroadcaster = websocket.NewBroadcaster()
	ProjectBroadcaster = websocket.NewBroadcaster()
}

type OpenFileRequest struct {
	ProjectName string `mapstructure:"project_name"`
	FilePath    string `mapstructure:"file_path"`
}

func OnOpenFile(ns *neffos.NSConn, body map[string]interface{}) error {
	req := OpenFileRequest{}
	err := mapstructure.Decode(body, &req)
	if err != nil {
		return err
	}
	FileBroadcaster.Subscribe(ns, util.FileKey(req.ProjectName, req.FilePath))

	notes := make([]model.Note, 0)
	result := model.DB.Where("project_name = ? AND file_path = ?", req.ProjectName, req.FilePath).Find(&notes)
	if result.Error != nil {
		return result.Error
	}
	for _, note := range notes {
		err = websocket.PushToClient(ns, "note", note)
		if err != nil {
			logrus.Warnf("WARN: Failed to push note to %s during open project", ns.String())
			continue
		}
	}

	return nil
}

type CloseFileRequest struct {
	ProjectName string `mapstructure:"project_name"`
	FilePath    string `mapstructure:"file_path"`
}

func OnCloseFile(ns *neffos.NSConn, body map[string]interface{}) error {
	req := CloseFileRequest{}
	err := mapstructure.Decode(body, &req)
	if err != nil {
		return err
	}
	FileBroadcaster.Unsubscribe(ns, util.FileKey(req.ProjectName, req.FilePath))
	return nil
}

type ProjectRequest struct {
	ProjectName string `mapstructure:"project_name"`
}

func OnOpenProject(ns *neffos.NSConn, body map[string]interface{}) error {
	req := ProjectRequest{}
	err := mapstructure.Decode(body, &req)
	if err != nil {
		return err
	}
	ProjectBroadcaster.Subscribe(ns, req.ProjectName)
	flows := make([]model.Flow, 0)
	result := model.DB.Where("project_name = ?", req.ProjectName).Find(&flows)
	if result.Error != nil {
		return result.Error
	}
	for _, flow := range flows {
		err = websocket.PushToClient(ns, "flow", flow)
		if err != nil {
			logrus.Warnf("WARN: Failed to push flow to %s during open project", ns.String())
			continue
		}
		rels := make([]model.FlowNoteRelation, 0)
		result = model.DB.Where("flow_id = ?", flow.ID).Find(&rels)
		if result.Error != nil {
			logrus.Warnf("Failed to get rels for flow %d", flow.ID)
			continue
		}
		for _, rel := range rels {
			err = websocket.PushToClient(ns, "flow_note_relation", rel)
			if err != nil {
				logrus.Warnf("Failed to push flow note relation to %s", ns.String())
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
		websocket.PushToClient(ns, "note", note)
	}
	return nil
}

func OnCloseProject(ns *neffos.NSConn, body map[string]interface{}) error {
	req := ProjectRequest{}
	err := mapstructure.Decode(body, &req)
	if err != nil {
		return err
	}
	ProjectBroadcaster.Unsubscribe(ns, req.ProjectName)
	return nil
}

func OnDisconnect(conn *neffos.Conn) {
	FileBroadcaster.Disconnect(conn.ID())
	ProjectBroadcaster.Disconnect(conn.ID())
}
