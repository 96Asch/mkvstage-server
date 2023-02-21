package service

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
)

type roleService struct {
	rr  domain.RoleRepository
	ur  domain.UserRepository
	urr domain.UserRoleRepository
}

func NewRoleService(rr domain.RoleRepository, ur domain.UserRepository, urr domain.UserRoleRepository) *roleService {
	return &roleService{
		rr:  rr,
		ur:  ur,
		urr: urr,
	}
}

func (rs roleService) FetchByID(ctx context.Context, id int64) (*domain.Role, error) {
	return nil, nil
}

func (rs roleService) FetchAll(ctx context.Context) (*[]domain.Role, error) {
	return nil, nil

}

func (rs roleService) Update(ctx context.Context, domain *domain.Role, principal *domain.User) error {
	return nil

}

func (rs roleService) UpdateBatch(ctx context.Context, domain *domain.Role, principal *domain.User) error {
	return nil

}

func (rs roleService) Store(ctx context.Context, domain *domain.Role, principal *domain.User) error {
	return nil

}

func (rs roleService) Remove(ctx context.Context, id int64, principal *domain.User) error {
	return nil

}
