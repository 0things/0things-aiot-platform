package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"

	"github.com/google/uuid"
)

var ErrInvalidRuleChainGraph = errors.New("rule chain graph must be a JSON object with nodes and edges")
var ErrRuleChainNameExists = errors.New("rule chain name already exists")

type RuleChainServiceInterface interface {
	List(ctx context.Context, page, size int, search, status string) ([]model.RuleChain, int64, error)
	Get(ctx context.Context, uuid string) (*model.RuleChain, error)
	Create(ctx context.Context, chain *model.RuleChain) error
	Update(ctx context.Context, chain *model.RuleChain) error
	Delete(ctx context.Context, uuid string) error
}

type RuleChainService struct {
	repo *repository.RuleChainRepository
}

func NewRuleChainService(repo *repository.RuleChainRepository) *RuleChainService {
	return &RuleChainService{repo: repo}
}

func (s *RuleChainService) List(ctx context.Context, page, size int, search, status string) ([]model.RuleChain, int64, error) {
	return s.repo.List(ctx, page, size, search, status)
}
func (s *RuleChainService) Get(ctx context.Context, id string) (*model.RuleChain, error) {
	return s.repo.Find(ctx, id)
}
func (s *RuleChainService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *RuleChainService) Create(ctx context.Context, chain *model.RuleChain) error {
	nameExists, err := s.repo.NameExists(ctx, chain.Name, "")
	if err != nil {
		return err
	}
	if nameExists {
		return ErrRuleChainNameExists
	}
	if err := validateRuleChainGraph(chain.Graph); err != nil {
		return err
	}
	if chain.UUID == "" {
		chain.UUID = uuid.NewString()
	}
	chain.Status = "draft"
	chain.Version = 1
	return s.repo.Create(ctx, chain)
}

func (s *RuleChainService) Update(ctx context.Context, chain *model.RuleChain) error {
	nameExists, err := s.repo.NameExists(ctx, chain.Name, chain.UUID)
	if err != nil {
		return err
	}
	if nameExists {
		return ErrRuleChainNameExists
	}
	if err := validateRuleChainGraph(chain.Graph); err != nil {
		return err
	}
	if chain.Version < 1 {
		chain.Version = 1
	}
	chain.Version++
	chain.Status = "draft"
	return s.repo.Save(ctx, chain)
}

func validateRuleChainGraph(raw json.RawMessage) error {
	if len(raw) == 0 || !json.Valid(raw) {
		return ErrInvalidRuleChainGraph
	}
	var graph struct {
		Nodes json.RawMessage `json:"nodes"`
		Edges json.RawMessage `json:"edges"`
	}
	if err := json.Unmarshal(raw, &graph); err != nil || len(graph.Nodes) == 0 {
		return ErrInvalidRuleChainGraph
	}
	var nodes []json.RawMessage
	if err := json.Unmarshal(graph.Nodes, &nodes); err != nil || len(nodes) == 0 {
		return ErrInvalidRuleChainGraph
	}
	return nil
}

func NormalizeRuleChainStatus(status string) string {
	status = strings.ToLower(strings.TrimSpace(status))
	if status == "published" {
		return status
	}
	return "draft"
}
