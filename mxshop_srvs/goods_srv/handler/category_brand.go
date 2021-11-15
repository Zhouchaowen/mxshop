package handler

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_srvs/goods_srv/proto"
)

//品牌分类
func (s *GoodsServer) CategoryBrandList(ctx context.Context, in *proto.CategoryBrandFilterRequest) (*proto.CategoryBrandListResponse, error){
	return nil, nil
}
//通过category获取brands
func (s *GoodsServer) GetCategoryBrandList(ctx context.Context, in *proto.CategoryInfoRequest) (*proto.BrandListResponse, error){
	return nil, nil
}
func (s *GoodsServer) CreateCategoryBrand(ctx context.Context, in *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error){
	return nil, nil
}
func (s *GoodsServer) DeleteCategoryBrand(ctx context.Context, in *proto.CategoryBrandRequest) (*emptypb.Empty, error){
	return nil, nil
}
func (s *GoodsServer) UpdateCategoryBrand(ctx context.Context, in *proto.CategoryBrandRequest) (*emptypb.Empty, error){
	return nil, nil
}
