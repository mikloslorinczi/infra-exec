package cmd

import (
	"github.com/mikloslorinczi/infra-exec/common"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func claimTask(task common.Task) (common.ResponseMsg, error) {
	var (
		responseMsg  common.ResponseMsg
		responseJSON []byte
		taskJSON     []byte
	)
	if err := common.ToJSON(&task, &taskJSON); err != nil {
		return responseMsg, errors.Wrap(err, "Cannot update task")
	}
	responseJSON, err := common.SendRequest("POST", viper.GetString("apiURL")+"/tasks/claim", taskJSON)
	if err != nil {
		return responseMsg, errors.Wrap(err, "Cannot update task")
	}
	if err := common.FromJSON(&responseMsg, responseJSON); err != nil {
		return responseMsg, errors.Wrap(err, "Cannot update task")
	}
	return responseMsg, nil
}
