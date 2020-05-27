package talert

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const version = "0.0.7"

type alerter struct {
	token  string
	chatID int

	client *http.Client
}

var (
	endpoint = "https://api.telegram.org/bot"
)

// SetEndpoint changes default endpoint. It's useful if you use proxy for avoid censorship
func SetEndpoint(e string) {
	endpoint = e
}

func (a *alerter) Alert(message string, fncs ...fieldFn) {
	url := fmt.Sprintf("%s%s/sendMessage?chat_id=%d&parse_mode=Markdown&text=%s",
		endpoint,
		a.token,
		a.chatID,
		render(message, fncs...))
	resp, err := a.client.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("read body: err=%s body=%s", err, bts)
		return
	}
	if err != nil {
		log.Printf("body: %s", bts)
	}
}

var defaultAlerter *alerter

func ParseDSN(dsn string) (string, int, error) {
	arr := strings.SplitN(dsn, "|", 2)
	if len(arr) != 2 {
		return "", 0, fmt.Errorf("bad dsn")
	}
	id, err := strconv.Atoi(arr[1])
	if err != nil {
		return "", 0, err
	}
	return arr[0], id, nil
}

func Init(token string, chatID int) error {
	defaultAlerter = &alerter{
		token:  token,
		chatID: chatID,
		client: &http.Client{Timeout: 10 * time.Second},
	}
	return nil
}

func Alert(message string, fncs ...fieldFn) {
	if defaultAlerter != nil {
		defaultAlerter.Alert(message, fncs...)
	}
}

type field struct {
	Name  string
	Value string
}

type fieldFn func() field

func String(k, v string) fieldFn {
	return func() field {
		return field{
			Name:  k,
			Value: v,
		}
	}
}

func Int(k string, v int) fieldFn {
	return func() field {
		return field{
			Name:  k,
			Value: strconv.Itoa(v),
		}
	}
}

func Error(k string, err error) fieldFn {
	var s string
	if err != nil {
		s = err.Error()
	}
	return func() field {
		return field{
			Name:  k,
			Value: s,
		}
	}
}

func Field(k string, v interface{}) fieldFn {
	return func() field {
		return field{
			Name:  k,
			Value: fmt.Sprint(v),
		}
	}
}

func render(message string, fncs ...fieldFn) string {
	result := message
	for i, fn := range fncs {
		if i == 0 {
			result += "\n\n"
		}
		f := fn()
		result += fmt.Sprintf("*%s*: %s\n", f.Name, f.Value)
	}
	return url.QueryEscape(result)
}
