package handlers

import (
	"context"
	"github.com/linger1216/go-utils/config"
	"github.com/linger1216/go-utils/log"

	pb "github.com/linger1216/jelly-doc/src/server/pb"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService(logger *log.Log, reader config.Reader) pb.ApiServer {
	return apiService{}
}

type apiService struct{}

func (s apiService) Create(ctx context.Context, in *pb.CreateApiRequest) (*pb.CreateApiResponse, error) {
	var resp pb.CreateApiResponse
	return &resp, nil
}

func (s apiService) Get(ctx context.Context, in *pb.GetApiRequest) (*pb.GetApiResponse, error) {
	var resp pb.GetApiResponse
	return &resp, nil
}

func (s apiService) List(ctx context.Context, in *pb.ListApiRequest) (*pb.ListApiResponse, error) {
	var resp pb.ListApiResponse
	return &resp, nil
}

func (s apiService) Update(ctx context.Context, in *pb.UpdateApiRequest) (*pb.EmptyResponse, error) {
	var resp pb.EmptyResponse
	return &resp, nil
}

func (s apiService) Delete(ctx context.Context, in *pb.DeleteApiRequest) (*pb.EmptyResponse, error) {
	var resp pb.EmptyResponse
	return &resp, nil
}
