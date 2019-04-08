package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

type SumoLogicHook struct {
	Url        string
	HttpClient *http.Client
	AppName    string
}

func NewSumo(config Config) Client {
	var client Client
	host, _ := os.Hostname()
	client.Logger = logrus.New()
	client.Logger.Formatter = &logrus.TextFormatter{
		ForceColors: false,
	}
	hook, err := NewSumoHook(config.Host, host)
	if err != nil {
		fmt.Println(err.Error())
		return client
	}
	client.Logger.Hooks.Add(hook)
	return client
}

func NewSumoHook(url string, appname string) (*SumoLogicHook, error) {
	if url == "" {
		return nil, fmt.Errorf("Unable to send logs to Sumo Logic. SUMO_ENDPOINT not provided")
	}
	client := &http.Client{}
	return &SumoLogicHook{url, client, appname}, nil
}

func (hook *SumoLogicHook) Fire(entry *logrus.Entry) error {
	data := make(logrus.Fields, len(entry.Data))
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}

	data["message"] = strings.Replace(entry.Message, "\"", "'", -1)
	data["level"] = entry.Level.String()
	s, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Failed to build json: %v", err)
	}
	err = hook.httpPost(s)
	if err != nil {
		return err
	}
	return nil
}

func (hook *SumoLogicHook) httpPost(s []byte) error {
	// already printed error about sumo_endpoint so be silent
	if hook.Url == "" || len(s) == 0 {
		// avoid panic and return if no url
		return nil
	}

	body := bytes.NewBuffer(s)
	req, err := http.NewRequest("POST", hook.Url, body)
	if err != nil {
		return fmt.Errorf("Error creating the request: %s", err.Error())
	}

	req.Close = true
	req.Header.Add("X-Sumo-Name", hook.AppName)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to post data: %s", err.Error())
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Failed to post data: %s", resp.Status)
	}
	return nil
}

func (s *SumoLogicHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}
