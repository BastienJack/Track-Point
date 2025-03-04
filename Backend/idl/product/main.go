package main

import (
	product "commerce/idl/product/kitex_gen/product/productservice"
	"log"
)

func main() {
	svr := product.NewServer(new(ProductServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
