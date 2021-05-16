package notion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"time"
)

const (
	apiBaseUrl    = "https://api.notion.com/"
	notionVersion = "2021-05-13"
)

type Client struct {
	AuthToken  string
	HTTPClient *http.Client
}

func (c Client) RetrievePage(pageId string) (*Page, error) {
	requestPath := path.Join("v1/pages/" + pageId)
	var page Page
	err := doNotionApi(c, requestPath, "GET", nil, &page)
	if err != nil {
		return nil, err
	}

	return &page, nil
}

func (c Client) RetrieveDatabase(databaseId string) (*Database, error) {
	requestPath := path.Join("v1/databases/" + databaseId)
	var database Database
	err := doNotionApi(c, requestPath, "GET", nil, &database)
	if err != nil {
		return nil, err
	}

	return &database, nil
}

// TODO
//func (c Client) RetrieveBlockChildren(blockId string) (*BlockList, error) {
//	requestPath := path.Join("v1/blocks/", blockId, "/children")
//	var blockList BlockList
//
//	// TODO: currently ignoring NextCursor even if it exists
//	err := doNotionApi(c, requestPath, "GET", nil, &blockList)
//	if err != nil {
//		return nil, err
//	}
//
//	return &blockList, nil
//}

func doNotionApi(c Client, path string, method string, requestData interface{}, result interface{}) error {
	uri := apiBaseUrl + path
	var jsonObj []byte
	var err error
	if requestData != nil {
		jsonObj, err = json.Marshal(requestData)
		if err != nil {
			return err
		}
	}
	body := bytes.NewBuffer(jsonObj)

	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return err
	}

	req.Header.Set("Notion-Version", notionVersion)
	if c.AuthToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.AuthToken))
	}
	var rsp *http.Response
	httpClient := c.getHTTPClient()
	rsp, err = httpClient.Do(req)
	if err != nil {
		return err
	}

	var d []byte
	d, _ = ioutil.ReadAll(rsp.Body)
	if rsp.StatusCode != 200 {
		return fmt.Errorf("Error: status code %s\nBody:\n%s\n", rsp.Status, d)
	}
	err = json.Unmarshal(d, result)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) getHTTPClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	httpClient := *http.DefaultClient
	httpClient.Timeout = time.Second * 30
	return &httpClient
}
