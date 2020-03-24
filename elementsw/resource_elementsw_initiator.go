package elementsw

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/fatih/structs"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/netapp/terraform-provider-netapp-elementsw/elementsw/element/jsonrpc"
)

// CreateInitiatorsRequest the user input for creating an initiator
type CreateInitiatorsRequest struct {
	Initiators []initiator `structs:"initiators"`
}

// CreateInitiatorsResult the API resutls for creatin an initiator
type CreateInitiatorsResult struct {
	Initiators []initiatorResponse `json:"initiators"`
}

// DeleteInitiatorsRequest the users input for deleteing an initiator
type DeleteInitiatorsRequest struct {
	Initiators []int `structs:"initiators"`
}

// ModifyInitiatorsRequest the users input for modifying a Request
type ModifyInitiatorsRequest struct {
	Initiators []initiator `structs:"initiators"`
}

func resourceElementSwInitiator() *schema.Resource {
	return &schema.Resource{
		Create: resourceElementSwInitiatorCreate,
		Read:   resourceElementSwInitiatorRead,
		Update: resourceElementSwInitiatorUpdate,
		Delete: resourceElementSwInitiatorDelete,
		Exists: resourceElementSwInitiatorExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"alias": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attributes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"volume_access_group_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"iqns": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceElementSwInitiatorCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating initiator: %#v", d)
	client := meta.(*Client)

	initiators := CreateInitiatorsRequest{}
	newInitiator := make([]initiator, 1)
	var iqns []string

	if v, ok := d.GetOk("name"); ok {
		newInitiator[0].Name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	if v, ok := d.GetOk("alias"); ok {
		newInitiator[0].Alias = v.(string)
	}

	if v, ok := d.GetOk("volume_access_group_id"); ok {
		newInitiator[0].VolumeAccessGroupID = v.(int)
	}

	if v, ok := d.GetOk("iqns"); ok {

		if a, ok := v.([]interface{}); ok {
			for i := range a {
				iqns = append(iqns, a[i].(string))
			}
		}
	}

	initiators.Initiators = newInitiator

	resp, err := createInitiators(client, initiators)
	if err != nil {
		log.Print("Error creating initiator")
		return err
	}

	d.SetId(fmt.Sprintf("%v", resp.Initiators[0].ID))
	log.Printf("Created initiator: %v %v", newInitiator[0].Name, resp.Initiators[0].ID)

	return resourceElementSwInitiatorRead(d, meta)
}

func createInitiators(client *Client, request CreateInitiatorsRequest) (CreateInitiatorsResult, error) {
	params := structs.Map(request)

	log.Printf("Parameters: %v", params)

	response, err := client.CallAPIMethod("CreateInitiators", params)
	if err != nil {
		log.Print("CreateInitiators request failed")
		return CreateInitiatorsResult{}, err
	}

	var result CreateInitiatorsResult
	if err := json.Unmarshal([]byte(*response), &result); err != nil {
		log.Print("Failed to unmarshall resposne from CreateInitiators")
		return CreateInitiatorsResult{}, err
	}
	log.Printf("Initiator Result: %v", result)
	return result, nil
}

func resourceElementSwInitiatorRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading initiator: %#v", d)
	client := meta.(*Client)

	initiators := listInitiatorRequest{}

	id := d.Id()
	s := make([]int, 1)
	convID, convErr := strconv.Atoi(id)

	if convErr != nil {
		return fmt.Errorf("id argument is required")
	}

	s[0] = convID
	initiators.Initiators = s

	res, err := client.listInitiators(initiators)
	if err != nil {
		return err
	}

	if len(res.Initiators) != 1 {
		return fmt.Errorf("Expected one Initiator to be found. Response contained %v results", len(res.Initiators))
	}

	d.Set("name", res.Initiators[0].Name)
	d.Set("alias", res.Initiators[0].Alias)
	d.Set("attributes", res.Initiators[0].Attributes)

	if len(res.Initiators[0].VolumeAccessGroups) == 1 {
		d.Set("volume_access_group_id", res.Initiators[0].VolumeAccessGroups[0])
	}

	return nil
}

func resourceElementSwInitiatorUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Updating initiator: %#v", d)
	client := meta.(*Client)

	initiators := ModifyInitiatorsRequest{}
	initiator := make([]initiator, 1)

	id := d.Id()
	convID, convErr := strconv.Atoi(id)

	if convErr != nil {
		return fmt.Errorf("id argument is required")
	}

	initiator[0].InitiatorID = convID

	if v, ok := d.GetOk("alias"); ok {
		initiator[0].Alias = v.(string)
	}

	if v, ok := d.GetOk("volume_access_group_id"); ok {
		initiator[0].VolumeAccessGroupID = v.(int)
	}

	initiators.Initiators = initiator

	err := modifyInitiators(client, initiators)
	if err != nil {
		return err
	}

	return nil
}

func modifyInitiators(client *Client, request ModifyInitiatorsRequest) error {
	params := structs.Map(request)

	_, err := client.CallAPIMethod("ModifyInitiators", params)
	if err != nil {
		log.Print("ModifyInitiators request failed")
		return err
	}

	return nil
}

func resourceElementSwInitiatorDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting initiator: %#v", d)
	client := meta.(*Client)

	initiators := DeleteInitiatorsRequest{}

	id := d.Id()
	s := make([]int, 1)
	convID, convErr := strconv.Atoi(id)

	if convErr != nil {
		return fmt.Errorf("id argument is required")
	}

	s[0] = convID
	initiators.Initiators = s

	err := deleteInitiator(client, initiators)
	if err != nil {
		return err
	}

	return nil
}

func deleteInitiator(client *Client, request DeleteInitiatorsRequest) error {
	params := structs.Map(request)

	_, err := client.CallAPIMethod("DeleteInitiators", params)
	if err != nil {
		log.Print("DeleteInitiator request failed")
		return err
	}

	return nil
}

func resourceElementSwInitiatorExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of initiator: %#v", d)
	client := meta.(*Client)

	initiators := listInitiatorRequest{}

	id := d.Id()
	s := make([]int, 1)
	convID, convErr := strconv.Atoi(id)

	if convErr != nil {
		return false, fmt.Errorf("id argument is required")
	}

	s[0] = convID
	initiators.Initiators = s

	res, err := client.listInitiators(initiators)
	if err != nil {
		if err, ok := err.(*jsonrpc.ResponseError); ok {
			if err.Name == "xUnknown" {
				d.SetId("")
				return false, nil
			}
		}
		return false, err
	}

	if len(res.Initiators) != 1 {
		d.SetId("")
		return false, nil
	}

	return true, nil
}
