package main

import (
	"flag"

	"github.com/kelseyhightower/envconfig"
)

func main() {
	config := struct {
		Port int `required:"true"`
	}{}

	flag.IntVar(&config.Port, "port", 0, "Set port via flag or via PORT=8080")
	flag.Parse()

	// port := flag.Int("port", 0, "Set port via flag or via PORT=8080")
	// flag.Parse()

	// if *port != 0 {
	// 	os.Setenv("PORT", strconv.Itoa(*port))
	// }

	err := envconfig.Process("", &config)
	if err != nil {
		panic(err)
	}
}
