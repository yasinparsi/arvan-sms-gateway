package provider

import (
	"encoding/json"
	"errors"
	"os"
	"sms-dispatcher/logger"
	"sms-dispatcher/model"
)

type ProviderFunc func(msg model.SmsMessage) error

var prefixMap map[string]ProviderFunc

func LoadOperatorConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	rawMap := map[string]string{}
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}

	prefixMap = make(map[string]ProviderFunc)
	for prefix, provider := range rawMap {
		switch provider {
		case "mci":
			prefixMap[prefix] = SendToMCI
		case "mtn":
			prefixMap[prefix] = SendToMTN
		case "rightel":
			prefixMap[prefix] = SendToRightel
		default:
			logger.Log.Warnw("Unknown provider in config", "prefix", prefix, "provider", provider)
		}
	}

	logger.Log.Infow("Operator config loaded", "count", len(prefixMap))
	return nil
}

func RouteAndSend(msg model.SmsMessage) error {
	if len(msg.Phone) < 4 {
		return errors.New("invalid phone number")
	}
	prefix := msg.Phone[:4]
	sendFunc, ok := prefixMap[prefix]
	if !ok {
		return errors.New("no provider matched for prefix: " + prefix)
	}
	return sendFunc(msg)
}
