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

// CreateAccountRequest is the user input for creating an account
type CreateAccountRequest struct {
	Username        string      `structs:"username"`
	InitiatorSecret string      `structs:"initiatorSecret,omitempty"`
	TargetSecret    string      `structs:"targetSecret,omitempty"`
	Attributes      interface{} `structs:"attributes,omitempty"`
}

// CreateAccountResult is the api returned output
type CreateAccountResult struct {
	Account account `json:"account"`
}

// ModifyAccountRequest is the users input for modifying an account
type ModifyAccountRequest struct {
	AccountID       int         `structs:"accountID"`
	InitiatorSecret string      `structs:"initiatorSecret,omitempty"`
	TargetSecret    string      `structs:"targetSecret,omitempty"`
	Attributes      interface{} `structs:"attributes,omitempty"`
	Username        string      `structs:"username,omitempty"`
}

// RemoveAccountRequest is the users input for deleteing an account
type RemoveAccountRequest struct {
	AccountID int `structs:"accountID"`
}

func resourceElementSwAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceElementSwAccountCreate,
		Read:   resourceElementSwAccountRead,
		Update: resourceElementSwAccountUpdate,
		Delete: resourceElementSwAccountDelete,
		Exists: resourceElementSwAccountExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"initiator_secret": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"target_secret": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"attributes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceElementSwAccountCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating account: %#v", d)
	client := meta.(*Client)

	acct := CreateAccountRequest{}

	if v, ok := d.GetOk("username"); ok {
		acct.Username = v.(string)
	} else {
		return fmt.Errorf("username argument is required")
	}

	if v, ok := d.GetOk("initiator_secret"); ok {
		acct.InitiatorSecret = v.(string)
	}

	if v, ok := d.GetOk("target_secret"); ok {
		acct.TargetSecret = v.(string)
	}

	resp, err := createAccount(client, acct)
	if err != nil {
		log.Print("Error creating account")
		return err
	}

	d.SetId(fmt.Sprintf("%v", resp.Account.AccountID))

	log.Printf("Created account: %v %v", acct.Username, resp.Account.AccountID)

	return resourceElementSwAccountRead(d, meta)
}

func createAccount(client *Client, request CreateAccountRequest) (CreateAccountResult, error) {
	params := structs.Map(request)

	log.Printf("Parameters: %v", params)

	response, err := client.CallAPIMethod("AddAccount", params)
	if err != nil {
		log.Print("CreateAccount request failed")
		return CreateAccountResult{}, err
	}

	var result CreateAccountResult
	if err := json.Unmarshal([]byte(*response), &result); err != nil {
		log.Print("Failed to unmarshall response from CreateAccount")
		return CreateAccountResult{}, err
	}
	return result, nil
}

func resourceElementSwAccountRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading account: %#v", d)
	client := meta.(*Client)

	id := d.Id()
	convID, convErr := strconv.Atoi(id)

	if convErr != nil {
		return fmt.Errorf("id argument is required")
	}

	res, err := client.getAccountByID(convID)
	if err != nil {
		log.Print("GetAccountByID failed")
		return err
	}

	if _, ok := d.GetOk("username"); ok {
		d.Set("username", res.Username)
	}

	if _, ok := d.GetOk("initiator_secret"); ok {
		d.Set("initiator_secret", res.InitiatorSecret)
	}

	if _, ok := d.GetOk("target_secret"); ok {
		d.Set("target_secret", res.TargetSecret)
	}

	return nil
}

func resourceElementSwAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Updating account %#v", d)
	client := meta.(*Client)

	acct := ModifyAccountRequest{}

	id := d.Id()
	convID, convErr := strconv.Atoi(id)

	if convErr != nil {
		return fmt.Errorf("id argument is required")
	}
	acct.AccountID = convID

	if v, ok := d.GetOk("username"); ok {
		acct.Username = v.(string)
	}

	if v, ok := d.GetOk("initiator_secret"); ok {
		acct.InitiatorSecret = v.(string)
	}

	if v, ok := d.GetOk("target_secret"); ok {
		acct.TargetSecret = v.(string)
	}

	err := modifyAccount(client, acct)
	if err != nil {
		return err
	}

	return nil
}

func modifyAccount(client *Client, request ModifyAccountRequest) error {
	params := structs.Map(request)

	_, err := client.CallAPIMethod("ModifyAccount", params)
	if err != nil {
		log.Print("ModifyAccount request failed")
		return err
	}

	return nil
}

func resourceElementSwAccountDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting account: %#v", d)
	client := meta.(*Client)

	acct := RemoveAccountRequest{}

	id := d.Id()
	convID, convErr := strconv.Atoi(id)

	if convErr != nil {
		return fmt.Errorf("id argument is required")
	}
	acct.AccountID = convID

	err := removeAccount(client, acct)
	if err != nil {
		return err
	}

	return nil
}

func removeAccount(client *Client, request RemoveAccountRequest) error {
	params := structs.Map(request)

	_, err := client.CallAPIMethod("RemoveAccount", params)
	if err != nil {
		log.Print("DeleteAccount request failed")
		return err
	}

	return nil
}

func resourceElementSwAccountExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of account: %#v", d)
	client := meta.(*Client)

	id := d.Id()
	convID, convErr := strconv.Atoi(id)

	if convErr != nil {
		return false, fmt.Errorf("id argument is required")
	}

	_, err := client.getAccountByID(convID)
	if err != nil {
		if err, ok := err.(*jsonrpc.ResponseError); ok {
			if err.Name == "xUnknownAccount" {
				d.SetId("")
				return false, nil
			}
		}
		log.Print("AccountExists failed")
		return false, err
	}

	return true, nil
}
