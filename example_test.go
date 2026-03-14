package envstruct_test

import (
	"fmt"

	"github.com/agentine/envstruct"
)

func ExampleProcess() {
	type Config struct {
		Host string
		Port int
	}
	var c Config
	err := envstruct.Process("APP", &c)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("ok")
	// Output: ok
}
