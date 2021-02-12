package handlers

import (
	"context"

	pb "github.com/linger1216/jelly-doc/src/server/pb"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.MemberServer {
	return memberService{}
}

type memberService struct{}

func (s memberService) Create(ctx context.Context, in *pb.CreateMemberRequest) (*pb.CreateMemberResponse, error) {
	var resp pb.CreateMemberResponse
	return &resp, nil
}

func (s memberService) Get(ctx context.Context, in *pb.GetMemberRequest) (*pb.GetMemberResponse, error) {
	var resp pb.GetMemberResponse
	return &resp, nil
}

func (s memberService) List(ctx context.Context, in *pb.ListMemberRequest) (*pb.ListMemberResponse, error) {
	var resp pb.ListMemberResponse
	return &resp, nil
}

func (s memberService) Update(ctx context.Context, in *pb.UpdateMemberRequest) (*pb.EmptyResponse, error) {
	var resp pb.EmptyResponse
	return &resp, nil
}

func (s memberService) Delete(ctx context.Context, in *pb.DeleteMemberRequest) (*pb.EmptyResponse, error) {
	var resp pb.EmptyResponse
	return &resp, nil
}
