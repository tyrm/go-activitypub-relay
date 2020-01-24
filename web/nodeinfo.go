package web

import (
	"encoding/json"
	"net/http"
	"time"
)

type NodeInfo struct {
	OpenRegistrations bool             `json:"openRegistrations"`
	Protocols         []string         `json:"protocols"`
	Services          NodeInfoServices `json:"services"`
	Software          NodeInfoSoftware `json:"software"`
	Usage             NodeInfoUsage    `json:"usage"`
	Version           string           `json:"version"`
	Metadata          NodeInfoMetadata `json:"metadata"`
}

type NodeInfoServices struct {
	Inbound  []string `json:"inbound"`
	Outbound []string `json:"outbound"`
}

type NodeInfoSoftware struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type NodeInfoUsage struct {
	LocalPosts int                `json:"localPosts"`
	Users      NodeInfoUsageUsers `json:"users"`
}

type NodeInfoUsageUsers struct {
	Total int `json:"total"`
}

type NodeInfoMetadata struct {
	Peers []string `json:"peers"`
}

func HandleNodeInfo(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	nodeInfo := NodeInfo{
		OpenRegistrations: true,
		Protocols:         []string{"activitypub"},
		Services: NodeInfoServices{
			Inbound:  []string{},
			Outbound: []string{},
		},
		Software: NodeInfoSoftware{
			Name:    "activityrelay",
			Version: "alpha+PettingZoo",
		},
		Usage: NodeInfoUsage{
			LocalPosts: 0,
			Users: NodeInfoUsageUsers{
				Total: 1,
			},
		},
		Version: "2.0",
		Metadata: NodeInfoMetadata{
			Peers: []string{},
		},
	}

	js, err := json.Marshal(nodeInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	elapsed := time.Since(start)
	logger.Infof("REQUEST HandleNodeInfo () %s", elapsed)
}
