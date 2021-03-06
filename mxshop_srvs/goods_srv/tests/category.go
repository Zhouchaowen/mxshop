package main

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestGetCategoryList() {
	rsp, err := brandClient.GetAllCategorysList(context.Background(), &emptypb.Empty{})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _, brand := range rsp.Data {
		fmt.Println(brand)
	}
}
