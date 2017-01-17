package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

type SumoLogicHook struct {
	Url             string
	HttpClient      *http.Client
	PendingMessages [][]byte
	AppName         string
}

func NewSumo(config Config) Client {
	var client Client
	host, _ := os.Hostname()
	client.Logger = logrus.New()
	client.Logger.Formatter = &logrus.TextFormatter{
		ForceColors: false,
	}
	hook, _ := NewSumoHook(config.Host, host)
	client.Logger.Hooks.Add(hook)

	return client
}

func NewSumoHook(url string, appname string) (*SumoLogicHook, error) {
	client := &http.Client{}
	return &SumoLogicHook{url, client, make([][]byte, 0), appname}, nil
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
	data["tstamp"] = entry.Time.Format(logrus.DefaultTimestampFormat)
	data["message"] = strings.Replace(entry.Message, "\"", "'", -1)
	data["level"] = entry.Level.String()

	s, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Failed to build json: %v", err)
	}
	// attempt to process pending messages first
	if len(hook.PendingMessages) != 0 {
		for i, m := range hook.PendingMessages {
			err := hook.httpPost(m)
			if err == nil {
				hook.PendingMessages, hook.PendingMessages[len(hook.PendingMessages)-1] = append(hook.PendingMessages[:i], hook.PendingMessages[i+1:]...), nil
			}
		}
	}
	err = hook.httpPost(s)
	if err != nil {
		// stash messages for next run
		hook.PendingMessages = append(hook.PendingMessages, s)
		return err
	}
	return nil
}

func (hook *SumoLogicHook) httpPost(s []byte) error {
	body := bytes.NewBuffer(s)
	req, err := http.NewRequest("POST", hook.Url, body)
	client := http.Client{}
	if req == nil {
		return fmt.Errorf("Something went wrong")
	}
	req.Header.Add("X-Sumo-Name", hook.AppName)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil || resp == nil {
		return fmt.Errorf("Failed to post data: %s", err.Error())
	} else if resp.StatusCode != 200 {
		return fmt.Errorf("Failed to post data: %s", resp.Status)
	} else {
		return nil
	}

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
