package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/hiroki-it/notify-slack-of-amplify-events/cmd/domain/detail"
	"github.com/hiroki-it/notify-slack-of-amplify-events/cmd/infrastructure/logger"
	"github.com/hiroki-it/notify-slack-of-amplify-events/cmd/usecase/services/amplify"
	"github.com/hiroki-it/notify-slack-of-amplify-events/cmd/usecase/services/notification"
)

// HandleRequest イベントをハンドリングします．
func HandleRequest(eventBridge events.CloudWatchEvent) (string, error) {

	log := logger.NewLogger()

	d := detail.NewDetail()

	// eventbridgeから転送されたJSONを構造体にマッピングします．
	err := json.Unmarshal([]byte(eventBridge.Detail), d)

	if err != nil {
		log.Error(err.Error())
		return fmt.Sprint("Failed to handle request"), err
	}

	amplifyApi, err := amplify.NewAmplifyAPI(os.Getenv("AWS_AMPLIFY_REGION"))

	if err != nil {
		log.Error(err.Error())
		return fmt.Sprint("Failed to handle request"), err
	}

	ac := amplify.NewAmplifyClient(amplifyApi)

	gbo, err := ac.GetBranchFromAmplify(d)

	if err != nil {
		log.Error(err.Error())
		return fmt.Sprint("Failed to handle request"), err
	}

	m := notification.NewMessage(
		d,
		gbo.Branch,
	)

	sm := m.BuildSlackMessage()

	sc := notification.NewSlackClient(
		&http.Client{},
		"https://slack.com/api/chat.postMessage",
	)

	sn := notification.NewSlackNotification(
		sc,
		sm,
	)

	err = sn.PostMessage()

	if err != nil {
		log.Error(err.Error())
		return fmt.Sprint("Failed to handle request"), err
	}

	return fmt.Sprint("Succeed to handle request"), nil
}
