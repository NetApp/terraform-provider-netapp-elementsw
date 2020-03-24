package elementsw

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/fatih/structs"
)

type listInitiatorRequest struct {
	Initiators []int `structs:"initiators"`
}

type listInitiatorResult struct {
	Initiators []initiatorResponse `json:"initiators"`
}

type initiator struct {
	Name                string      `structs:"name,omitempty"`
	Alias               string      `structs:"alias,omitempty"`
	Attributes          interface{} `structs:"attributes,omitempty"`
	VolumeAccessGroupID int         `structs:"volumeAccessGroupID,omitempty"`
	InitiatorID         int         `structs:"initiatorID,omitempty"`
}

type initiatorResponse struct {
	Name               string      `json:"initiatorName"`
	Alias              string      `json:"alias"`
	Attributes         interface{} `json:"attributes"`
	ID                 int         `json:"initiatorID"`
	VolumeAccessGroups []int       `json:"volumeAccessGroups"`
}

func (c *Client) getInitiatorByID(id string) (initiator, error) {
	convID, err := strconv.Atoi(id)
	if err != nil {
		return initiator{}, err
	}

	initID := make([]int, 1)
	initID[0] = convID

	params := structs.Map(listInitiatorRequest{Initiators: initID})

	response, err := c.CallAPIMethod("ListInitiators", params)
	if err != nil {
		log.Print("ListInitiators request failed")
		return initiator{}, err
	}

	var result listInitiatorResult
	if err := json.Unmarshal([]byte(*response), &result); err != nil {
		log.Print("Failed to unmarshal response from ListInitiators")
		return initiator{}, err
	}

	if len(result.Initiators) != 1 {
		return initiator{}, fmt.Errorf("Expected one Initiator to be found. Response contained %v results", len(result.Initiators))
	}

	var initiator initiator
	initiator.Name = result.Initiators[0].Name
	initiator.Alias = result.Initiators[0].Alias
	initiator.Attributes = result.Initiators[0].Attributes
	initiator.InitiatorID = result.Initiators[0].ID
	if len(result.Initiators[0].VolumeAccessGroups) == 1 {
		initiator.VolumeAccessGroupID = result.Initiators[0].VolumeAccessGroups[0]
	}

	return initiator, nil
}

func (c *Client) listInitiators(request listInitiatorRequest) (listInitiatorResult, error) {
	params := structs.Map(request)

	response, err := c.CallAPIMethod("ListInitiators", params)
	if err != nil {
		log.Print("ListInitiators request failed")
		return listInitiatorResult{}, err
	}

	var result listInitiatorResult
	if err := json.Unmarshal([]byte(*response), &result); err != nil {
		log.Print("Failed to unmarshall response from ListInitiators")
		return listInitiatorResult{}, err
	}

	return result, nil
}
