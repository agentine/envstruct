package envstruct_test

import (
	"fmt"
	"os"

	"github.com/agentine/envstruct"
)

func ExampleProcess() {
	if err := os.Setenv("APP_HOST", "localhost"); err != nil {
		panic(err)
	}
	if err := os.Setenv("APP_PORT", "8080"); err != nil {
		panic(err)
	}
	defer func() { _ = os.Unsetenv("APP_HOST") }()
	defer func() { _ = os.Unsetenv("APP_PORT") }()

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
	fmt.Printf("Host=%s Port=%d\n", c.Host, c.Port)
	// Output: Host=localhost Port=8080
}

func ExampleProcess_nested() {
	if err := os.Setenv("APP_DATABASE_HOST", "db.local"); err != nil {
		panic(err)
	}
	if err := os.Setenv("APP_DATABASE_PORT", "5432"); err != nil {
		panic(err)
	}
	defer func() { _ = os.Unsetenv("APP_DATABASE_HOST") }()
	defer func() { _ = os.Unsetenv("APP_DATABASE_PORT") }()

	type DB struct {
		Host string
		Port int
	}
	type Config struct {
		Database DB
	}
	var c Config
	err := envstruct.Process("APP", &c)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("DB=%s:%d\n", c.Database.Host, c.Database.Port)
	// Output: DB=db.local:5432
}

func ExampleProcess_required() {
	type Config struct {
		Secret string `env:"SECRET,required"`
	}
	var c Config
	err := envstruct.Process("APP", &c)
	if err != nil {
		fmt.Println("got expected error")
		return
	}
	fmt.Println("unexpected success")
	// Output: got expected error
}

func ExampleUsage() {
	type Config struct {
		Host  string `default:"localhost" desc:"Server hostname"`
		Port  int    `default:"8080" desc:"Server port"`
		Debug bool   `desc:"Enable debug mode"`
	}
	_ = envstruct.Usage("APP", &Config{}, os.Stdout)
	// Output:
	//   KEY        TYPE    DEFAULT               DESCRIPTION
	//   APP_HOST   string  [default: localhost]  Server hostname
	//   APP_PORT   int     [default: 8080]       Server port
	//   APP_DEBUG  bool                          Enable debug mode
}
