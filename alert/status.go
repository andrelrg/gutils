package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/astrolink/gutils/general"
)

const (
	StandardStatusWebhook         = "standard"
	SubscriptionStatusWebhook     = "subscription"
	I18nSubscriptionStatusWebhook = "i18nSubscription"
	I18nNotificationStatusWebhook = "i18nNotification"
	PushNotificationStatusWebhook = "pushNotification"
)

type StatusAlert struct {
	WebhookList map[string]interface{}
}

func getSupportedChannels() []string {
	return []string{StandardStatusWebhook,
		SubscriptionStatusWebhook,
		I18nSubscriptionStatusWebhook,
		I18nNotificationStatusWebhook,
		PushNotificationStatusWebhook}
}

// NewStatusAlert makes a new instance of StatusAlert struct
func NewStatusAlert() *StatusAlert {
	list := make(map[string]interface{})
	return &StatusAlert{WebhookList: list}
}

// SetWebhook associate a channel webhook name to a channel webhook url
func (s *StatusAlert) SetWebhook(name, url string) error {
	var err error

	if in, _ := general.InArray(name, getSupportedChannels()); !in {
		err = fmt.Errorf("%s channel is not in supported channel webhook list, supported: [%s]", name, strings.Join(getSupportedChannels(), ","))
		log.Println(err)
		return err
	}

	s.WebhookList[name] = url
	return err
}

// SendSlackMessage sends a message to a Slack channel through a webhook
func (s *StatusAlert) SendSlackMessage(messageContent, webhookName string) {
	var err error
	var endpoint string

	if in, _ := general.InArray(webhookName, getSupportedChannels()); !in {
		err = fmt.Errorf("%s channel is not in supported channel webhookName list, supported: [%s]", webhookName, strings.Join(getSupportedChannels(), ","))
		log.Println(err)
		return
	}

	endpoint = s.WebhookList[webhookName].(string)

	if len(endpoint) == 0 {
		log.Printf("empty endpoint for webhookName %s\n", webhookName)
		return
	}

	binName := os.Args[0]
	binName = filepath.Base(binName)
	textMessage := make(map[string]interface{})
	textMessage["text"] = fmt.Sprintf("%s - %s", binName, messageContent)

	body, err := json.Marshal(textMessage)

	if err != nil {
		err = fmt.Errorf("error parsing merssage to json")
		log.Println(err)
		return
	}

	client := http.Client{}
	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))

	if err != nil {
		err = fmt.Errorf("error creating request, " + err.Error())
		log.Println(err)
		return
	}

	response, err := client.Do(request)

	if err != nil {
		err = fmt.Errorf("error making request, " + err.Error())
		log.Println(err)
		return
	}

	_, err = ioutil.ReadAll(response.Body)

	if err != nil {
		err = fmt.Errorf("error reading response body, " + err.Error())
		log.Println(err)
		return
	}
}
