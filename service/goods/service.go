package goods

import (
	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/adaptor/repo/goods"
)

type Service struct {
	lesson goods.ILesson
}

func NewService(adaptor adaptor.IAdaptor) *Service {
	return &Service{
		lesson: goods.NewLesson(adaptor),
	}
}
