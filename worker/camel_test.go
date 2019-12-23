package worker

import (
	"fmt"
	"testing"
)

func TestCamelIt(t *testing.T) {
	fmt.Println(CamelIt("t__user_info___xx", true))
	fmt.Println(CamelIt("t__user_info___xx", false))
}
