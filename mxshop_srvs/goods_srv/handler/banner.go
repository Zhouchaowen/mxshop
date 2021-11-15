package handler

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_srvs/goods_srv/proto"
)

//轮播图
func (s *GoodsServer) BannerList(ctx context.Context, in *emptypb.Empty) (*proto.BannerListResponse, error){
	return nil, nil
}
func (s *GoodsServer) CreateBanner(ctx context.Context, in *proto.BannerRequest) (*proto.BannerResponse, error){
	return nil, nil
}
func (s *GoodsServer) DeleteBanner(ctx context.Context, in *proto.BannerRequest) (*emptypb.Empty, error){
	return nil, nil
}
func (s *GoodsServer) UpdateBanner(ctx context.Context, in *proto.BannerRequest) (*emptypb.Empty, error){
	return nil, nil
}
