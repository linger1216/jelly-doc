package handlers

import (
	"context"
	"github.com/linger1216/go-utils/code"
	"github.com/linger1216/go-utils/config"
	"github.com/linger1216/go-utils/db/postgres"
	"github.com/linger1216/go-utils/log"
	"github.com/linger1216/jelly-doc/src/server/core"

	pb "github.com/linger1216/jelly-doc/src/server/pb"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService(logger *log.Log, reader config.Reader) pb.BasisApiServer {
	uri := reader.GetString("postgres", "uri")
	db := postgres.NewPostgres(postgres.NewConfig(uri))
	return &BasisApiService{core.NewApiDBService(logger, db)}
}

type BasisApiService struct {
	dbProxy pb.BasisApiServer
}

func (s *BasisApiService) Create(ctx context.Context, in *pb.CreateApiRequest) (*pb.CreateApiResponse, error) {
	if in == nil || in.Apis == nil || len(in.Apis) == 0 {
		return nil, code.ErrInvalidPara
	}
	return s.dbProxy.Create(ctx, in)
}

func (s *BasisApiService) Get(ctx context.Context, in *pb.GetApiRequest) (*pb.GetApiResponse, error) {
	if in == nil || len(in.Ids) == 0 {
		return nil, code.ErrInvalidPara
	}
	return s.dbProxy.Get(ctx, in)
}

func (s *BasisApiService) List(ctx context.Context, in *pb.ListApiRequest) (*pb.ListApiResponse, error) {
	if in == nil {
		return nil, code.ErrInvalidPara
	}

	if in.CurrentPage < 0 {
		in.CurrentPage = 0
	}

	if in.PageSize <= 0 {
		in.PageSize = 10
	}

	return s.dbProxy.List(ctx, in)
}

func (s *BasisApiService) Update(ctx context.Context, in *pb.UpdateApiRequest) (*pb.EmptyResponse, error) {
	if in == nil || in.Apis == nil || len(in.Apis) == 0 {
		return nil, code.ErrInvalidPara
	}
	return s.dbProxy.Update(ctx, in)
}

func (s *BasisApiService) Delete(ctx context.Context, in *pb.DeleteApiRequest) (*pb.EmptyResponse, error) {
	if in == nil || len(in.Ids) == 0 {
		return nil, code.ErrInvalidPara
	}
	return s.dbProxy.Delete(ctx, in)
}
