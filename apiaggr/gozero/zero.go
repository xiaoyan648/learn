package gozero

import (
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/mr"
)

func HomePage() {
	var userInfo *User
	var productList []Product
	_ = mr.Finish(
		func() (err error) {
			userInfo, err = getUser()
			return err
		},
		func() (err error) {
			productList, err = getProductList()
			return err
		},
	)
	fmt.Printf("用户信息:%+v\n", userInfo)
	fmt.Printf("商品信息:%+v\n", productList)
}

/********用户服务**********/

type User struct {
	Name string
	Age  uint8
}

func getUser() (*User, error) {
	time.Sleep(500 * time.Millisecond)
	var u User
	u.Name = "wuqinqiang"
	u.Age = 18
	return &u, nil
}

/********商品服务**********/

type Product struct {
	Title string
	Price uint32
}

func getProductList() ([]Product, error) {
	time.Sleep(400 * time.Millisecond)
	var list []Product
	list = append(list, Product{
		Title: "SHib",
		Price: 10,
	})
	return list, nil
}
