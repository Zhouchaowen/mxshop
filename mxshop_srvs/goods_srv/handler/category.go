package handler

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_srvs/goods_srv/proto"
)

//商品分类
func (s *GoodsServer) GetAllCategorysList(ctx context.Context, in *emptypb.Empty) (*proto.CategoryListResponse, error){
	return nil, nil
}
//获取子分类
func (s *GoodsServer) GetSubCategory(ctx context.Context, in *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error){
	return nil, nil
}
func (s *GoodsServer) CreateCategory(ctx context.Context, in *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error){
	return nil, nil
}
func (s *GoodsServer) DeleteCategory(ctx context.Context, in *proto.DeleteCategoryRequest) (*emptypb.Empty, error){
	return nil, nil
}
func (s *GoodsServer) UpdateCategory(ctx context.Context, in *proto.CategoryInfoRequest) (*emptypb.Empty, error){
	return nil, nil
}
