package elementsw

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/fatih/structs"
)

type getAccountByIDRequest struct {
	AccountID int `structs:"accountID"`
}

type getAccountByNameRequest struct {
	AccountName string `structs:"username"`
}

type getAccountResult struct {
	Account account `json:"account"`
}

type account struct {
	AccountID       int         `json:"accountID"`
	Attributes      interface{} `json:"attributes"`
	InitiatorSecret string      `json:"initiatorSecret"`
	Status          string      `json:"status"`
	TargetSecret    string      `json:"targetSecret"`
	Username        string      `json:"username"`
}

func (c *Client) getAccountByID(id int) (account, error) {
	params := structs.Map(getAccountByIDRequest{AccountID: id})
	return c.getAccountDetails(params, "GetAccountByID")
}

func (c *Client) getAccountByName(name string) (account, error) {
	params := structs.Map(getAccountByNameRequest{AccountName: name})
	return c.getAccountDetails(params, "GetAccountByName")
}

func (c *Client) getAccountDetails(params map[string]interface{}, method string) (account, error) {
	response, err := c.CallAPIMethod(method, params)
	if err != nil {
		log.Print(method + " request failed")
		return account{}, err
	}

	var result getAccountResult
	if err := json.Unmarshal([]byte(*response), &result); err != nil {
		log.Print("Failed to unmarshal response from GetAccountByID")
		return account{}, err
	}

	return result.Account, nil
}

func (c *Client) getAccountByIDOrName(account interface{}) (account, error) {
	if id, err := strconv.Atoi(account.(string)); err == nil {
		// when account is numeric
		accountDetails, err := c.getAccountByID(id)
		if err == nil {
			return accountDetails, nil
		}
		// when account name is numeric
		accountDetails, err = c.getAccountByName(account.(string))
		if err == nil {
			return accountDetails, nil
		}
		return accountDetails, err
	}
	accountDetails, err := c.getAccountByName(account.(string))
	if err == nil {
		return accountDetails, nil
	}
	return accountDetails, err
}
