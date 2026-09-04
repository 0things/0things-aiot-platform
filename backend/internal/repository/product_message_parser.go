package repository

import (
	"context"
	"errors"

	"aiot-backend/internal/model"
	"aiot-backend/internal/tenant"

	"gorm.io/gorm"
)

type ProductMessageParserRepository struct {
	db *gorm.DB
}

func NewProductMessageParserRepository(db *gorm.DB) *ProductMessageParserRepository {
	return &ProductMessageParserRepository{db: db}
}

func (r *ProductMessageParserRepository) FindByProductID(ctx context.Context, productID int64) (*model.ProductMessageParser, error) {
	q := useQuery(r.db)
	parser, err := q.ProductMessageParser.WithContext(ctx).Where(
		q.ProductMessageParser.ProductID.Eq(productID),
		q.ProductMessageParser.OrganizationID.Eq(tenant.GetOrganizationID(ctx)),
	).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return parser, nil
}

func (r *ProductMessageParserRepository) Save(ctx context.Context, parser *model.ProductMessageParser) error {
	parser.OrganizationID = tenant.GetOrganizationID(ctx)
	if existing, err := r.FindByProductID(ctx, parser.ProductID); err == nil {
		parser.ID = existing.ID
		return useQuery(r.db).ProductMessageParser.WithContext(ctx).Save(parser)
	} else if !errors.Is(err, ErrNotFound) {
		return err
	}
	return useQuery(r.db).ProductMessageParser.WithContext(ctx).Create(parser)
}
