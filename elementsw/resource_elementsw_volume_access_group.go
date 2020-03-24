package elementsw

import (
	"fmt"
	"log"
	"strconv"

	"encoding/json"

	"github.com/fatih/structs"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/netapp/terraform-provider-netapp-elementsw/elementsw/element/jsonrpc"
)

// CreateVolumeAccessGroupRequest the users input for creating an volume access group
type CreateVolumeAccessGroupRequest struct {
	Name       string      `structs:"name"`
	Initiators []string    `structs:"initiators"`
	Volumes    []int       `structs:"volumes"`
	Attributes interface{} `structs:"attributes"`
	ID         int         `structs:"id"`
}

// CreateVolumeAccessGroupResult the API results for creating an aaccess group
type CreateVolumeAccessGroupResult struct {
	VolumeAccessGroupID int `json:"volumeAccessGroupID"`
	volumeAccessGroup
}

// DeleteVolumeAccessGroupRequest the user input for deleteing an volume access group
type DeleteVolumeAccessGroupRequest struct {
	VolumeAccessGroupID    int  `structs:"volumeAccessGroupID"`
	DeleteOrphanInitiators bool `structs:"deleteOrphanInitiators"`
	Force                  bool `structs:"force"`
}

// ModifyVolumeAccessGroupRequest the users input for modifying a colume access group
type ModifyVolumeAccessGroupRequest struct {
	VolumeAccessGroupID    int         `structs:"volumeAccessGroupID"`
	Name                   string      `structs:"name"`
	Attributes             interface{} `structs:"attributes"`
	Initiators             []int       `structs:"initiators"`
	DeleteOrphanInitiators bool        `structs:"deleteOrphanInitiators"`
	Volumes                []int       `structs:"volumes"`
}

func resourceElementSwVolumeAccessGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceElementSwVolumeAccessGroupCreate,
		Read:   resourceElementSwVolumeAccessGroupRead,
		Update: resourceElementSwVolumeAccessGroupUpdate,
		Delete: resourceElementSwVolumeAccessGroupDelete,
		Exists: resourceElementSwVolumeAccessGroupExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"volumes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"attributes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"initiators": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceElementSwVolumeAccessGroupCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating volume access group: %#v", d)
	client := meta.(*Client)

	vag := CreateVolumeAccessGroupRequest{}

	if v, ok := d.GetOk("name"); ok {
		vag.Name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	if raw, ok := d.GetOk("volumes"); ok {
		for _, v := range raw.([]interface{}) {
			vag.Volumes = append(vag.Volumes, v.(int))
		}
	}

	resp, err := createVolumeAccessGroup(client, vag)
	if err != nil {
		log.Print("Error creating volume access group")
		return err
	}

	d.SetId(fmt.Sprintf("%v", resp.VolumeAccessGroupID))
	log.Printf("Created volume access group: %v %v", vag.Name, resp.VolumeAccessGroupID)

	return resourceElementSwVolumeAccessGroupRead(d, meta)
}

func createVolumeAccessGroup(client *Client, request CreateVolumeAccessGroupRequest) (CreateVolumeAccessGroupResult, error) {
	params := structs.Map(request)

	log.Printf("Parameters: %v", params)

	response, err := client.CallAPIMethod("CreateVolumeAccessGroup", params)
	if err != nil {
		log.Print("CreateVolumeAccessGroup request failed")
		return CreateVolumeAccessGroupResult{}, err
	}

	var result CreateVolumeAccessGroupResult
	if err := json.Unmarshal([]byte(*response), &result); err != nil {
		log.Print("Failed to unmarshall response from CreateVolumeAccessGroup")
		return CreateVolumeAccessGroupResult{}, err
	}
	return result, nil
}

func resourceElementSwVolumeAccessGroupRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading volume access group: %#v", d)
	client := meta.(*Client)

	vags := listVolumeAccessGroupsRequest{}

	id := d.Id()
	s := make([]int, 1)
	convID, convErr := strconv.Atoi(id)

	if convErr != nil {
		return fmt.Errorf("id argument is required")
	}

	s[0] = convID
	vags.VolumeAccessGroups = s

	res, err := client.listVolumeAccessGroups(vags)
	if err != nil {
		return err
	}

	if len(res.VolumeAccessGroupsNotFound) > 0 {
		return fmt.Errorf("Unable to find Volume Access Groups with the ID of %v", res.VolumeAccessGroupsNotFound)
	}

	if len(res.VolumeAccessGroups) != 1 {
		return fmt.Errorf("Expected one Volume Access Group to be found. Response contained %v results", len(res.VolumeAccessGroups))
	}

	d.Set("name", res.VolumeAccessGroups[0].Name)
	d.Set("initiators", res.VolumeAccessGroups[0].Initiators)
	d.Set("volumes", res.VolumeAccessGroups[0].Volumes)

	return nil
}

func resourceElementSwVolumeAccessGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Updating volume access group %#v", d)
	client := meta.(*Client)

	vag := ModifyVolumeAccessGroupRequest{}

	id := d.Id()
	convID, convErr := strconv.Atoi(id)

	if convErr != nil {
		return fmt.Errorf("id argument is required")
	}
	vag.VolumeAccessGroupID = convID

	if v, ok := d.GetOk("name"); ok {
		vag.Name = v.(string)

	} else {
		return fmt.Errorf("name argument is required during update")
	}

	if raw, ok := d.GetOk("volumes"); ok {
		for _, v := range raw.([]interface{}) {
			vag.Volumes = append(vag.Volumes, v.(int))
		}
	} else {
		return fmt.Errorf("expecting an array of volume ids to change")
	}

	err := modifyVolumeAccessGroup(client, vag)
	if err != nil {
		return err
	}

	return nil
}

func modifyVolumeAccessGroup(client *Client, request ModifyVolumeAccessGroupRequest) error {
	params := structs.Map(request)

	_, err := client.CallAPIMethod("ModifyVolumeAccessGroup", params)
	if err != nil {
		log.Print("ModifyVolumeAccessGroup request failed")
		return err
	}

	return nil
}

func resourceElementSwVolumeAccessGroupDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting volume access group: %#v", d)
	client := meta.(*Client)

	vag := DeleteVolumeAccessGroupRequest{}

	id := d.Id()
	convID, convErr := strconv.Atoi(id)

	if convErr != nil {
		return fmt.Errorf("id argument is required")
	}
	vag.VolumeAccessGroupID = convID

	err := deleteVolumeAccessGroup(client, vag)
	if err != nil {
		return err
	}

	return nil
}

func deleteVolumeAccessGroup(client *Client, request DeleteVolumeAccessGroupRequest) error {
	params := structs.Map(request)

	_, err := client.CallAPIMethod("DeleteVolumeAccessGroup", params)
	if err != nil {
		log.Print("DeleteVolumeAccessGroup request failed")
		return err
	}

	return nil
}

func resourceElementSwVolumeAccessGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of volume access group: %#v", d)
	client := meta.(*Client)

	vags := listVolumeAccessGroupsRequest{}

	id := d.Id()
	s := make([]int, 1)
	convID, convErr := strconv.Atoi(id)

	if convErr != nil {
		return false, fmt.Errorf("id argument is required")
	}

	s[0] = convID
	vags.VolumeAccessGroups = s

	res, err := client.listVolumeAccessGroups(vags)
	if err != nil {
		if err, ok := err.(*jsonrpc.ResponseError); ok {
			if err.Name == "xUnknown" {
				d.SetId("")
				return false, nil
			}
		}
		return false, err
	}

	if len(res.VolumeAccessGroupsNotFound) > 0 {
		d.SetId("")
		return false, nil
	}

	if len(res.VolumeAccessGroups) != 1 {
		d.SetId("")
		return false, nil
	}

	return true, nil
}
