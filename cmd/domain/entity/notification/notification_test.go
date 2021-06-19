package notification

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hiroki-it/notify-slack-of-amplify-events/cmd/domain/entity/event"
	"github.com/stretchr/testify/assert"

	aws_amplify "github.com/aws/aws-sdk-go/service/amplify"
)

/**
 * ステータスがSUCCEED場合に，PostMessageメソッドが成功することをテストします．
 */
func TestPostMessage(t *testing.T) {

	t.Helper()

	// テストケース
	cases := []struct {
		// テストケース名
		name string
		// 期待値
		expected error
		// テストデータ
		eventDetail     *event.EventDetail
		getBranchOutput *aws_amplify.GetBranchOutput
	}{
		{
			name:     "TestPostMessage_JobStatusSucceed_ReturnNil",
			expected: nil,
			eventDetail: &event.EventDetail{
				AppId:      "1",
				BranchName: "test",
				JobId:      "1",
				JobStatus:  "SUCCEED",
			},
			getBranchOutput: &aws_amplify.GetBranchOutput{
				Branch: &aws_amplify.Branch{DisplayName: aws.String("feature-test")},
			},
		},
		{
			name:     "TestPostMessage_JobStatusFailed_ReturnNil",
			expected: nil,
			eventDetail: &event.EventDetail{
				AppId:      "1",
				BranchName: "test",
				JobId:      "1",
				JobStatus:  "FAILED",
			},
			getBranchOutput: &aws_amplify.GetBranchOutput{
				Branch: &aws_amplify.Branch{DisplayName: aws.String("feature-test")},
			},
		},
	}

	// 外部サーバのモックを構築します．
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "{}")
	}))

	defer ts.Close()

	// 反復処理で全てのテストケースを検証します．
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			event := event.NewEvent(tt.eventDetail)

			message := NewMessage(
				event,
				tt.getBranchOutput.Branch,
			)

			slackClient := NewSlackClient(
				&http.Client{},
				ts.URL, // モックサーバのURLに差し替えます．
			)

			slackMessage := message.BuildSlackMessage()

			slackNotification := NewSlackNotification(
				slackClient,
				slackMessage,
			)

			err := slackNotification.PostMessage()

			if err != nil {
				t.Fatal(err.Error())
			}

			assert.Nil(t, err)
		})
	}
}