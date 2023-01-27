package cmd

import (
	"flag"
	"fmt"
	"lib2snip/entities"
	"log"

	"github.com/spf13/viper"
)

func Run() {
	// ...
	flag.Parse()
	fmt.Println(flag.Lookup("config").Value.String())
	configFile := flag.Lookup("config").Value.String()
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file : %s\n", configFile)
		// panic(err)
	}
	fmt.Println(viper.AllKeys())
	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
	snippet := entities.Snippet{}
	viper.Unmarshal(&snippet)
	fmt.Println(snippet)
}
func init() {

}
func init() {
	flag.String("config", "config.yaml", "config file")
	flag.Bool("debug", false, "debug mode")
	flag.String("output", "output.json", "output file")
}
