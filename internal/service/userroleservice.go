package service

import (
	"context"

	"github.com/96Asch/mkvstage-server/internal/domain"
)

type userRoleService struct {
	urr domain.UserRoleRepository
	rr  domain.RoleRepository
}

func NewUserRoleService(urr domain.UserRoleRepository, rr domain.RoleRepository) *userRoleService {
	return &userRoleService{
		urr: urr,
		rr:  rr,
	}
}

func (urs userRoleService) UpdateBatch(ctx context.Context, userRoles *[]domain.UserRole, principal *domain.User) error {
	return nil
}
