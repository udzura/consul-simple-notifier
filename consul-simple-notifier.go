package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pelletier/go-toml"
	//"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

type config struct {
	emails     []string
	ikachanUrl string
	channel    string
}

type consulAlert struct {
	Timestamp string
	Node      string
	ServiceId string
	Service   string
	CheckId   string
	Check     string
	Status    string
	Output    string
	Notes     string
}

func (c *consulAlert) TrimmedOutput() string {
	return strings.TrimSpace(c.Output)
}

const (
	version         = "0.0.1"
	ircBodyTemplate = setIrcMode(ircBold) +
		"{{.Service}}({{.CheckId}}) is now {{.StatusString}}" +
		setIrcMode(ircBold) +
		"on {{.NodeString}}" +
		" - {{.TrimmedOutput}}"
	mailTitleTemplate = "Check {{.CheckId}} is now {{.Status}} on {{.Node}}"
	mailBodyTemplate  = `
{{.Service}}({{.CheckId}}) is now {{.Status}}
On node {{.Node}}

Output is:
  {{.TrimmedOutput}}
`
)

var (
	logger = log.New(os.Stdout, "[consul-simple-notifier] ", log.LstdFlags)
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
	conf.channel = parsed.Get("ikachan.channel").(string)
	logger.Printf("conf is: %+v\n", conf)

	err = json.NewDecoder(os.Stdin).Decode(&input)
	if err != nil {
		panic(err.Error())
	}
	logger.Printf("input json is: %+v\n", input)

	for _, content := range input {
		notifyEmail(conf.emails, content)
		notifyIkachan(conf.ikachanUrl, conf.channel, content)
	}
}

func notifyEmail(recipients []string, content consulAlert) error {
	for _, address := range recipients {
		var titleBuf, bodyBuf bytes.Buffer
		titleTmpl := template.Must(template.New("emailTitle").Parse(mailTitleTemplate))
		bodyTmpl := template.Must(template.New("emailBody").Parse(mailBodyTemplate))
		err := titleTmpl.Execute(&titleBuf, &content)
		err = bodyTmpl.Execute(&bodyBuf, &content)
		if err != nil {
			return err
		}
		title := titleBuf.String()

		logger.Printf("Sending... %s to %s\n", title, address)
		cmd := exec.Command("/bin/mail", "-s", title, address)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return err
		}

		if err := cmd.Start(); err != nil {
			return err
		}

		fmt.Fprint(stdin, bodyBuf.String())
		stdin.Close()
		logger.Printf("Send!\n")
		cmd.Wait()
	}
	return nil
}

func notifyIkachan(ikachanUrl string, channel string, content consulAlert) error {
	joinUrl := fmt.Sprintf("%s/join", ikachanUrl)
	noticeUrl := fmt.Sprintf("%s/notice", ikachanUrl)

	values := make(url.Values)
	values.Set("channel", channel)

	resp1, err := http.PostForm(joinUrl, values)
	defer resp1.Body.Close()
	if err != nil {
		return err
	}

	var bodyBuf bytes.Buffer
	bodyTmpl := template.Must(template.New("ircBody").Parse(ircBodyTemplate))
	err = bodyTmpl.Execute(&bodyBuf, &content)
	if err != nil {
		return err
	}
	body := bodyBuf.String()

	values.Set("message", body)

	logger.Printf("Posted! %+v", values)
	resp2, err := http.PostForm(noticeUrl, values)
	defer resp2.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

func showVersion() {
	fmt.Printf("consul-simple-notifier version: %s\n", version)
}
