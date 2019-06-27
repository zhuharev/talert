package talert

import (
	"fmt"
	"net/http"
	"strconv"
)

const version = "0.0.1"

type alerter struct {
	token  string
	chatID int
}

const (
	endpoint = "https://api.telegram.org/bot"
)

func (a *alerter) Alert(message string, fncs ...fieldFn) {
	url := fmt.Sprintf("%s%s/sendMessage?chat_id=%d&parse_mode=Markdown&text=%s",
		endpoint,
		a.token,
		a.chatID,
		render(message, fncs...))
	_, err := http.Get(url)
	if err != nil {
		_ = err
	}
}

var defaultAlerter *alerter

func Init(token string, chatID int) error {
	defaultAlerter = &alerter{
		token:  token,
		chatID: chatID,
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
		result += fmt.Sprintf("**%s**: %s\n", f.Name, f.Value)
	}
	return result
}
