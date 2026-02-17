package provider

import (
	"encoding/json"
	"fmt"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

func getByOwnerName[T any](c *casdoorsdk.Client, action, owner, name string) (*T, error) {
	queryMap := map[string]string{
		"id": fmt.Sprintf("%s/%s", owner, name),
	}
	url := c.GetUrl(action, queryMap)
	bytes, err := c.DoGetBytes(url)
	if err != nil {
		return nil, err
	}
	var result *T
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
