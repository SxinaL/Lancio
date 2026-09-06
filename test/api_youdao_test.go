package test

import (
	"Lancio/api"
	"fmt"
	"testing"

)

func TestRequestYoudao(t *testing.T) {
	

	body, err := api.RequestYoudao("experiments","en")
	if err != nil {
	  	fmt.Printf("请求失败: %v", err)
	}
	fmt.Println(string(body))
	
}

