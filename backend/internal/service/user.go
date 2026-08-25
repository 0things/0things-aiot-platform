package service

import (
	"context"
	"strings"
	"time"

	v1 "aiot-backend/api/v1"
	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"
	"aiot-backend/internal/tenant"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, req *v1.RegisterRequest) error
	Login(ctx context.Context, req *v1.LoginRequest) (string, error)
	GetProfile(ctx context.Context, userId string) (*v1.GetProfileResponseData, error)
	UpdateProfile(ctx context.Context, userId string, req *v1.UpdateProfileRequest) error
	ListMyOrganizations(ctx context.Context, userId string) ([]*v1.OrganizationItem, error)
	SwitchOrganization(ctx context.Context, userId string, orgID int64) (string, error)
}

func NewUserService(
	service *Service,
	userRepo repository.UserRepository,
	orgRepo repository.OrganizationRepository,
	orgUserRepo repository.OrganizationUserRepository,
) UserService {
	return &userService{
		userRepo:    userRepo,
		orgRepo:     orgRepo,
		orgUserRepo: orgUserRepo,
		Service:     service,
	}
}

type userService struct {
	userRepo    repository.UserRepository
	orgRepo     repository.OrganizationRepository
	orgUserRepo repository.OrganizationUserRepository
	*Service
}

func defaultOrgName(email string) string {
	if idx := strings.Index(email, "@"); idx > 0 {
		return email[:idx] + "'s Org"
	}
	return email
}

func (s *userService) Register(ctx context.Context, req *v1.RegisterRequest) error {
	// check username
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return v1.ErrInternalServerError
	}
	if err == nil && user != nil {
		return v1.ErrEmailAlreadyUse
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// Generate user ID
	userId, err := s.sid.GenString()
	if err != nil {
		return err
	}
	user = &model.User{
		UserId:   userId,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	org := &model.Organization{
		Name: defaultOrgName(req.Email),
		Slug: userId,
	}

	return s.tm.Transaction(ctx, func(ctx context.Context) error {
		if err = s.userRepo.Create(ctx, user); err != nil {
			return err
		}
		if err = s.orgRepo.Create(ctx, org); err != nil {
			return err
		}
		orgUser := &model.OrganizationUser{
			OrganizationID: org.Id,
			UserID:         userId,
		}
		if err = s.orgUserRepo.Create(ctx, orgUser); err != nil {
			return err
		}
		return nil
	})
}

func (s *userService) Login(ctx context.Context, req *v1.LoginRequest) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return "", v1.ErrUnauthorized
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", v1.ErrUnauthorized
	}

	orgUsers, err := s.orgUserRepo.ListByUser(ctx, user.UserId)
	if err != nil {
		return "", err
	}

	var selectedOrgID int64
	now := time.Now()

	if len(orgUsers) == 0 {
		// 历史用户无归属时，新建个人组织并关联
		org := &model.Organization{
			Name: defaultOrgName(user.Email),
			Slug: user.UserId,
		}
		err = s.tm.Transaction(ctx, func(ctx context.Context) error {
			if err := s.orgRepo.Create(ctx, org); err != nil {
				return err
			}
			orgUser := &model.OrganizationUser{
				OrganizationID: org.Id,
				UserID:         user.UserId,
				LastLoginAt:    &now,
			}
			return s.orgUserRepo.Create(ctx, orgUser)
		})
		if err != nil {
			return "", err
		}
		selectedOrgID = org.Id
	} else {
		// 查找最近登录的组织，若全为 NULL 则取 min(org_id)
		var latestOrgUser *model.OrganizationUser
		var minOrgID int64 = orgUsers[0].OrganizationID
		var latestTime *time.Time

		for _, ou := range orgUsers {
			if ou.OrganizationID < minOrgID {
				minOrgID = ou.OrganizationID
			}
			if ou.LastLoginAt != nil {
				if latestTime == nil || ou.LastLoginAt.After(*latestTime) {
					latestTime = ou.LastLoginAt
					latestOrgUser = ou
				}
			}
		}

		if latestOrgUser != nil {
			selectedOrgID = latestOrgUser.OrganizationID
		} else {
			selectedOrgID = minOrgID
		}

		if err := s.orgUserRepo.UpdateLastLogin(ctx, user.UserId, selectedOrgID, now); err != nil {
			s.logger.WithContext(ctx).Warn("failed to update user last login time", zap.Error(err))
		}
	}

	token, err := s.jwt.GenToken(user.UserId, selectedOrgID, time.Now().Add(time.Hour*24*90))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) ListMyOrganizations(ctx context.Context, userId string) ([]*v1.OrganizationItem, error) {
	orgUsers, err := s.orgUserRepo.ListByUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	if len(orgUsers) == 0 {
		return []*v1.OrganizationItem{}, nil
	}

	orgIDs := make([]int64, 0, len(orgUsers))
	for _, ou := range orgUsers {
		orgIDs = append(orgIDs, ou.OrganizationID)
	}

	orgs, err := s.orgRepo.ListByIDs(ctx, orgIDs)
	if err != nil {
		return nil, err
	}

	currentOrgID := tenant.GetOrganizationID(ctx)
	items := make([]*v1.OrganizationItem, 0, len(orgs))
	for _, org := range orgs {
		items = append(items, &v1.OrganizationItem{
			Id:        org.Id,
			Name:      org.Name,
			IsCurrent: org.Id == currentOrgID,
		})
	}
	return items, nil
}

func (s *userService) SwitchOrganization(ctx context.Context, userId string, orgID int64) (string, error) {
	isMember, err := s.orgUserRepo.IsMember(ctx, userId, orgID)
	if err != nil {
		return "", err
	}
	if !isMember {
		return "", v1.ErrForbidden
	}

	now := time.Now()
	if err := s.orgUserRepo.UpdateLastLogin(ctx, userId, orgID, now); err != nil {
		s.logger.WithContext(ctx).Warn("failed to update user last login time on switch", zap.Error(err))
	}

	token, err := s.jwt.GenToken(userId, orgID, time.Now().Add(time.Hour*24*90))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *userService) GetProfile(ctx context.Context, userId string) (*v1.GetProfileResponseData, error) {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &v1.GetProfileResponseData{
		UserId:   user.UserId,
		Nickname: user.Nickname,
		Email:    user.Email,
	}, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userId string, req *v1.UpdateProfileRequest) error {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	user.Email = req.Email
	user.Nickname = req.Nickname

	if err = s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

