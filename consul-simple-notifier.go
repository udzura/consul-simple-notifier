package main

import (
	//"encoding/json"
	"flag"
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"os"
)

type config struct {
	emails     []string
	ikachanUrl string
}

type consulAlert struct {
	timestamp string
	node      string
	serviceId string
	service   string
	checkId   string
	check     string
	output    string
	notes     string
}

func main() {
	var (
		configPath string
		conf       config
		//input      []consulAlert
	)

	flag.StringVar(&configPath, "c", "/etc/consul-simple-notifier.ini", "Config path")
	flag.Parse()

	parsed, err := toml.LoadFile(configPath)
	if err != nil {
		panic(err.Error())
	}
	recipients := parsed.Get("email.recipients")
	for _, address := range recipients.([]interface{}) {
		conf.emails = append(conf.emails, address.(string))
	}
	conf.ikachanUrl = parsed.Get("ikachan.url").(string)
	fmt.Printf("%+v\n", conf)

	in, _ := ioutil.ReadAll(os.Stdin)
	fmt.Println(string(in))

	/*
		err = json.NewDecoder(os.Stdin).Decode(&input)
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("%+v\n", input)
	*/
}
