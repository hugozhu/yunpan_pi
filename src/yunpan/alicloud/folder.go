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

func (c *Client) MakeFolder(parentId int64, name string) (*Folder, error) {
	params := &url.Values{}
	params.Set("pid", fmt.Sprintf("%d", parentId))
	params.Set("name", name)

	result, err := c.PostCall("/folder/mkdir", params)
	if err != nil {
		return nil, err
	}
	var folder Folder
	json.Unmarshal(result, &folder)
	if !folder.Suc {
		return nil, ApiError{ErrorCode: 0, ErrorDescription: fmt.Sprintf("Failed to make folder: %d %s", parentId, name)}
	}
	return &folder, err
}

func (c *Client) RemoveFolder(id int64) (*Folder, error) {
	params := &url.Values{}
	params.Set("id", fmt.Sprintf("%d", id))

	result, err := c.PostCall("/folder/remove", params)
	if err != nil {
		return nil, err
	}
	var folder Folder
	json.Unmarshal(result, &folder)
	if !folder.Suc {
		return nil, ApiError{ErrorCode: 0, ErrorDescription: fmt.Sprintf("Failed to remove folder: %d", id)}
	}
	return &folder, err
}

func (c *Client) MoveFolder(id int64, parentId int64) (*Folder, error) {
	params := &url.Values{}
	params.Set("id", fmt.Sprintf("%d", id))
	params.Set("pid", fmt.Sprintf("%d", parentId))
	result, err := c.PostCall("/folder/move", params)
	if err != nil {
		return nil, err
	}
	var folder Folder
	json.Unmarshal(result, &folder)
	if !folder.Suc {
		return nil, ApiError{ErrorCode: 0, ErrorDescription: fmt.Sprintf("Failed to move folder: %d", id)}
	}
	return &folder, err
}

func (c *Client) RenameFolder(id int64, newname string) (*Folder, error) {
	params := &url.Values{}
	params.Set("id", fmt.Sprintf("%d", id))
	params.Set("newName", newname)
	result, err := c.PostCall("/folder/rename", params)
	if err != nil {
		return nil, err
	}
	var folder Folder
	json.Unmarshal(result, &folder)
	if !folder.Suc {
		return nil, ApiError{ErrorCode: 0, ErrorDescription: fmt.Sprintf("Failed to rename folder: %d", id)}
	}
	return &folder, err
}
