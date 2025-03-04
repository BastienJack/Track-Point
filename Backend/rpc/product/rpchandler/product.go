package rpchandler

import (
	"commerce/idl/product/kitex_gen/product"
	"context"
)

// ProductServiceImpl implements the last userservice interface defined in the IDL.
type ProductServiceImpl struct{}

// SearchProduct implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) SearchProduct(ctx context.Context, req *product.SearchProductRequest) (resp *product.SearchProductResponse, err error) {
	return
}

// GetProduct implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) GetProduct(ctx context.Context, req *product.GetProductRequest) (resp *product.GetProductResponse, err error) {
	// TODO: Your code here...
	return
}

// GetProductList implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) GetProductList(ctx context.Context, req *product.GetProductListRequest) (resp *product.GetProductListResponse, err error) {
	// TODO: Your code here...
	return
}
