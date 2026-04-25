package goods

import (
	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/adaptor/repo/admin"
	"github.com/JunLang-7/mall/adaptor/repo/goods"
	"github.com/JunLang-7/mall/adaptor/rpc"
)

type Service struct {
	lesson  goods.ILesson
	user    admin.IAdminUser
	storage rpc.IStorage
}

func NewService(adaptor adaptor.IAdaptor) *Service {
	return &Service{
		lesson:  goods.NewLesson(adaptor),
		user:    admin.NewRepo(adaptor),
		storage: rpc.NewStorage(adaptor),
	}
}
