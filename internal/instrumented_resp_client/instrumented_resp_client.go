package instrumented_resp_client

import (
	"strings"

	resp_value "github.com/codecrafters-io/redis-tester/internal/resp/value"
	resp_client "github.com/codecrafters-io/redis-tester/internal/resp_client"
	testerutils "github.com/codecrafters-io/tester-utils"
)

func NewInstrumentedRespClient(stageHarness *testerutils.StageHarness, addr string, clientIdentifier string) (*resp_client.RespClient, error) {
	logPrefix := ""
	if clientIdentifier != "" {
		logPrefix = clientIdentifier + ": "
	}

	return resp_client.NewRespClientWithCallbacks(
		addr,
		resp_client.RespClientCallbacks{
			OnSendCommand: func(command string, args ...string) {
				if len(args) > 0 {
					stageHarness.Logger.Infof("%s$ redis-cli %s %s", logPrefix, command, strings.Join(args, " "))
				} else {
					stageHarness.Logger.Infof("%s$ redis-cli %s", logPrefix, command)
				}
			},
			OnBytesSent: func(bytes []byte) {
				stageHarness.Logger.Debugf("%sSent bytes: %q", logPrefix, string(bytes))
			},
			OnBytesReceived: func(bytes []byte) {
				stageHarness.Logger.Debugf("%sReceived bytes: %q", logPrefix, string(bytes))
			},
			OnValueRead: func(value resp_value.Value) {
				stageHarness.Logger.Debugf("%sReceived RESP value: %s", logPrefix, value.FormattedString())
			},
		},
	)
}