package rest

import (
	"encoding/json"
	"fmt"
	"github.com/taekion-org/sawtooth-client-sdk-go/transport/types"
	"net/url"
)

// batchStatusRestResponse represents a REST API reply when batch status is requested.
type batchStatusRestResponse struct {
	Data []struct {
		Id                  string				`json:"id"`
		Status				types.BatchStatus	`json:"status"`
		InvalidTransactions []struct {
			Id		string						`json:"id"`
			Message		string					`json:"message"`
		} `json:"invalid_transactions"`

	} `json:"data"`
}

// GetBatchStatus returns the status for a single batch.
func (self *SawtoothClientTransportRest) GetBatchStatus(batchId string, wait int) (types.BatchStatus, error) {
	statusMap, err := self.GetBatchStatusMultiple([]string{batchId}, wait)
	if err != nil {
		return "", err
	}

	return statusMap[batchId], nil
}

// GetBatchStatusMultiple returns the statuses for a list of batches.
func (self *SawtoothClientTransportRest) GetBatchStatusMultiple(batchIds []string, wait int) (map[string]types.BatchStatus, error) {
	relativeUrl := &url.URL{Path: "/batch_statuses"}

	var waitParam string
	if wait == 0 {
		waitParam = "false"
	} else {
		waitParam = fmt.Sprint("%d", wait)
	}

	query := relativeUrl.Query()
	query.Add("wait", waitParam)
	relativeUrl.RawQuery = query.Encode()

	body, err := json.Marshal(batchIds)
	if err != nil {
		return nil, err
	}

	data, err := self.doPostRequestJson(relativeUrl, body)
	if err != nil {
		return nil, err
	}

	var response batchStatusRestResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	resultMap := make(map[string]types.BatchStatus, len(batchIds))
	for _, result := range response.Data {
		resultMap[result.Id] = result.Status
	}

	return resultMap, nil
}
