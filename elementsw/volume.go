package elementsw

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/fatih/structs"
)

type listVolumesRequest struct {
	Volumes               []int `structs:"volumeIDs"`
	IncludeVirtualVolumes bool  `structs:"includeVirtualVolumes"`
}

type listVolumesResult struct {
	Volumes []volume `json:"volumes"`
}

// listVolumesByVolumeNameRequest this is the list of volumes for given Account and volume name
type listVolumesByVolumeNameRequest struct {
	VolumeName string `structs:"volumeName"`
	Accounts   []int  `structs:"accounts"`
}

// listVolumesByAccountIDRequest this is the list of volumes for given Account and volume ID
type listVolumesByAccountIDRequest struct {
	Accounts []int `structs:"accounts"`
}

type volume struct {
	Name     string `json:"name"`
	VolumeID int    `json:"volumeID"`
	Iqn      string `json:"iqn"`
}

func (c *Client) listVolumesByVolumeID(request listVolumesByAccountIDRequest) (listVolumesResult, error) {
	params := structs.Map(request)
	return c.getVolumesDetails(params)
}

func (c *Client) listVolumesByVolumeName(request listVolumesByVolumeNameRequest) (listVolumesResult, error) {
	params := structs.Map(request)
	return c.getVolumesDetails(params)
}

func (c *Client) listVolumes(request listVolumesRequest) (listVolumesResult, error) {
	params := structs.Map(request)
	return c.getVolumesDetails(params)
}

func (c *Client) getVolumesDetails(params map[string]interface{}) (listVolumesResult, error) {

	response, err := c.CallAPIMethod("ListVolumes", params)
	if err != nil {
		log.Print("ListVolumes request failed")
		return listVolumesResult{}, err
	}

	var result listVolumesResult
	if err := json.Unmarshal([]byte(*response), &result); err != nil {
		log.Print("Failed to unmarshall response from ListVolumes")
		return listVolumesResult{}, err
	}

	return result, nil
}

func (c *Client) getVolumeByIDOrName(volumeInfo interface{}, accountID int) (listVolumesResult, error) {
	if id, err := strconv.Atoi(volumeInfo.(string)); err == nil {
		volumeInput := listVolumesByAccountIDRequest{}
		accIds := make([]int, 1)
		accIds[0] = accountID
		volumeInput.Accounts = accIds
		volumeDetails, err := c.listVolumesByVolumeID(volumeInput)
		if err == nil {
			if len(volumeDetails.Volumes) > 0 {
				for _, vol := range volumeDetails.Volumes {
					if vol.VolumeID == id {
						volDetail := listVolumesResult{}
						volDetail.Volumes = make([]volume, 1)
						volDetail.Volumes[0] = vol
						return volDetail, nil
					}
				}
			}
		}
		volumeInputName := listVolumesByVolumeNameRequest{}
		accIds = make([]int, 1)
		accIds[0] = accountID
		volumeInputName.Accounts = accIds
		volumeInputName.VolumeName = volumeInfo.(string)
		volumeDetailsForName, err := c.listVolumesByVolumeName(volumeInputName)
		if err == nil {
			if len(volumeDetailsForName.Volumes) != 1 {
				return listVolumesResult{}, fmt.Errorf("Expected one Volume to be found. Response contained %v results", len(volumeDetailsForName.Volumes))
			}
			return volumeDetailsForName, nil
		}
		return volumeDetailsForName, err
	}
	volumeInput := listVolumesByVolumeNameRequest{}
	accIds := make([]int, 1)
	accIds[0] = accountID
	volumeInput.Accounts = accIds
	volumeInput.VolumeName = volumeInfo.(string)
	volumeDetails, err := c.listVolumesByVolumeName(volumeInput)
	if err == nil {
		if len(volumeDetails.Volumes) != 1 {
			return listVolumesResult{}, fmt.Errorf("Expected one Volume to be found. Response contained %v results", len(volumeDetails.Volumes))
		}
		return volumeDetails, nil
	}
	return volumeDetails, err
}
