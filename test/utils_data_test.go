package test

import (
	"Lancio/utils"
	"fmt"
	"testing"
)
func TestIsValidVocabulary(t *testing.T) {
	// 测试有效词汇
	flag := utils.IsValidWord("in this section","en") 
	fmt.Println(flag)
}