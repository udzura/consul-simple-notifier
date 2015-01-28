package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pelletier/go-toml"
	//"io/ioutil"
	"os"
	"os/exec"
)

type config struct {
	emails     []string
	ikachanUrl string
}

type consulAlert struct {
	Timestamp string
	Node      string
	ServiceId string
	Service   string
	CheckId   string
	Check     string
	Output    string
	Notes     string
}

const (
	version = "0.0.1"
)

func main() {
	var (
		justShowVersion bool
		configPath      string
		conf            config
		input           []consulAlert
	)

	flag.BoolVar(&justShowVersion, "v", false, "Show version")
	flag.BoolVar(&justShowVersion, "version", false, "Show version")

	flag.StringVar(&configPath, "c", "/etc/consul-simple-notifier.ini", "Config path")
	flag.Parse()

	if justShowVersion {
		showVersion()
		return
	}

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

	err = json.NewDecoder(os.Stdin).Decode(&input)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v\n", input)

	for _, content := range input {
		notifyEmail(conf.emails, content)
	}
}

func notifyEmail(recipients []string, content consulAlert) error {
	for _, address := range recipients {
		fmt.Printf("Sending... %s to %+v\n", address, content)
		cmd := exec.Command("/bin/mail", "-s", "Alert from consul", address)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return err
		}

		if err := cmd.Start(); err != nil {
			return err
		}

		fmt.Fprint(stdin, "This is a sample mail\n")
		stdin.Close()
		fmt.Printf("Send!\n")
		cmd.Wait()
	}
	return nil
}

func notifyIkachan(ikachanUrl string, content consulAlert) {
}

func showVersion() {
	fmt.Printf("consul-simple-notifier version: %s\n", version)
}
