package converter_test

import (
	"fmt"
	"github.com/milin2436/BaiduPCS-Go/pcsutil/converter"
	"testing"
)

func TestShortDisplay(t *testing.T) {
	for i := 0; i < 20; i++ {
		fmt.Println([]byte(converter.ShortDisplay("\u0000我我\u0000我我我我", i)))
	}
}
