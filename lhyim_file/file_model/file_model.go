package file_model

import (
	"github.com/google/uuid"
	"lhyim_server/common/models"
)

type FileModel struct {
	models.Model
	Uid      uuid.UUID `json:"uid"` //文件唯一id /api/file/{uuid}
	UserID   uint      `json:"userID"`
	FileName string    `json:"fileName"`
	Size     int64     `json:"size"`
	Path     string    `json:"path"`
	Hash     string    `json:"hash"` //文件hash
}

func (file *FileModel) WebPath() string {
	return "/api/file/" + file.Uid.String()
}
