package models

type Node struct {
	NodeName   string `json:"node_name"`
	NodeId     string `json:"node_id"`
	NodeStatus int    `json:"node_status"`
}
