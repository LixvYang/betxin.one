package snowflake

import (
	"fmt"
	"testing"
)

func Test_Snowflake(t *testing.T) {
	Init("2023-11-10", 2)
	fmt.Println(GenID())
	fmt.Println(GenID())
	fmt.Println(GenID())
	fmt.Println(GenID())
	fmt.Println(GenID())
}
