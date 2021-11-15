package handler

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_srvs/goods_srv/proto"
)

//品牌和轮播图
func (s *GoodsServer) BrandList(ctx context.Context, in *proto.BrandFilterRequest) (*proto.BrandListResponse, error){
	return nil, nil
}
func (s *GoodsServer) CreateBrand(ctx context.Context, in *proto.BrandRequest) (*proto.BrandInfoResponse, error){
	return nil, nil
}
func (s *GoodsServer) DeleteBrand(ctx context.Context, in *proto.BrandRequest) (*emptypb.Empty, error){
	return nil, nil
}
func (s *GoodsServer) UpdateBrand(ctx context.Context, in *proto.BrandRequest) (*emptypb.Empty, error){
	return nil, nil
}
