package elementsw

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/fatih/structs"
)

type listVolumeAccessGroupsRequest struct {
	VolumeAccessGroups []int `structs:"volumeAccessGroups"`
}

type listVolumeAccessGroupsResult struct {
	VolumeAccessGroups         []volumeAccessGroup `json:"volumeAccessGroups"`
	VolumeAccessGroupsNotFound []int               `json:"volumeAccessGroupsNotFound"`
}

type volumeAccessGroup struct {
	VolumeAccessGroupID int      `json:"volumeAccessGroupID"`
	Name                string   `json:"name"`
	Initiators          []string `json:"initiators"`
	Volumes             []int    `json:"volumes"`
	ID                  int      `json:"id"`
}

func (c *Client) getVolumeAccessGroupByID(id string) (volumeAccessGroup, error) {
	convID, err := strconv.Atoi(id)
	if err != nil {
		return volumeAccessGroup{}, err
	}

	vagIDs := make([]int, 1)
	vagIDs[0] = convID

	params := structs.Map(listVolumeAccessGroupsRequest{VolumeAccessGroups: vagIDs})

	response, err := c.CallAPIMethod("ListVolumeAccessGroups", params)
	if err != nil {
		log.Print("ListVolumeAccessGroups request failed")
		return volumeAccessGroup{}, err
	}

	var result listVolumeAccessGroupsResult
	if err := json.Unmarshal([]byte(*response), &result); err != nil {
		log.Print("Failed to unmarshal respone from ListVolumeAccessGroups")
		return volumeAccessGroup{}, err
	}

	if len(result.VolumeAccessGroupsNotFound) > 0 {
		return volumeAccessGroup{}, fmt.Errorf("Unable to find Volume Access Groups with the ID of %v", result.VolumeAccessGroupsNotFound)
	}

	if len(result.VolumeAccessGroups) != 1 {
		return volumeAccessGroup{}, fmt.Errorf("Expected one Volume Access Group to be found. Response contained %v results", len(result.VolumeAccessGroups))
	}

	return result.VolumeAccessGroups[0], nil
}

func (c *Client) listVolumeAccessGroups(request listVolumeAccessGroupsRequest) (listVolumeAccessGroupsResult, error) {
	params := structs.Map(request)

	response, err := c.CallAPIMethod("ListVolumeAccessGroups", params)
	if err != nil {
		log.Print("ListVolumeAccessGroups request failed")
		return listVolumeAccessGroupsResult{}, err
	}

	var result listVolumeAccessGroupsResult
	if err := json.Unmarshal([]byte(*response), &result); err != nil {
		log.Print("Failed to unmarshall response from ListVolumeAccessGroups")
		return listVolumeAccessGroupsResult{}, err
	}

	return result, nil
}
