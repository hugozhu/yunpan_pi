package alicloud

import (
	"encoding/json"
	"fmt"
	"net/url"
)

func (c *Client) FolderList(dirId int64) (*FileList, error) {
	params := &url.Values{}
	params.Set("dirId", fmt.Sprintf("%d", dirId))
	result, err := c.GetCall("/folder/list", params)
	if err != nil {
		return nil, err
	}
	var fileList FileList
	json.Unmarshal(result, &fileList)
	return &fileList, err
}
