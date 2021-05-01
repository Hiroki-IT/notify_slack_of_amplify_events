package unit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	aws_amplify "github.com/aws/aws-sdk-go/service/amplify"
	"github.com/hiroki-it/notify-slack-of-amplify-events/cmd/entities/amplify"
	"github.com/hiroki-it/notify-slack-of-amplify-events/cmd/entities/eventbridge"
	"github.com/hiroki-it/notify-slack-of-amplify-events/cmd/entities/slack"
	m_amplify "github.com/hiroki-it/notify-slack-of-amplify-events/test/mock/amplify"
	"github.com/stretchr/testify/assert"
)

func SlackResponse(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "200")
}

/**
 * 関数をテストします．
 */
func TestLambdaHandler(t *testing.T) {

	input := aws_amplify.GetBranchInput{
		AppId:      aws.String("123456789"),
		BranchName: aws.String("feature/test"),
	}

	mockedAPI := new(m_amplify.MockedAmplifyAPI)

	// スタブに引数として渡される値と，その時の返却値を定義する．
	mockedAPI.On("GetBranch", &input).Return(Branch{DisplayName: aws.String("feature-test")}, nil)

	client := amplify.NewAmplifyClient(mockedAPI)

	eventDetail := &eventbridge.EventDetail{
		AppId:      "123456789",
		BranchName: "feature/test",
		JobId:      "123456789",
		JobStatus:  "SUCCEED",
	}

	// 検証対象の関数を実行する．スタブを含む一連の処理が実行される．
	response, _ := client.GetBranchFromAmplify(eventDetail)

	slackClient := slack.NewSlackClient()

	message := slackClient.BuildMessage(
		eventDetail,
		&slack.AmplifyBranch{DisplayName: aws.StringValue(response.Branch.DisplayName)},
	)

	json, _ := json.Marshal(message)

	request := httptest.NewRequest(
		"POST",
		"https://slack.com/api/chat.postMessage",
		bytes.NewBuffer(json),
	)

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("SLACK_API_TOKEN")))

	// HTTPリクエストを送信する．
	writer := httptest.NewRecorder()

	assert.Equal(t, http.StatusOK, writer)
}