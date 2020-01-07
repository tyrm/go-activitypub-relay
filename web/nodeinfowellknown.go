package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type NodeInfoWellKnown struct {
	Links []NodeInfoWellKnownLink `json:"links"`
}

type NodeInfoWellKnownLink struct {
	Rel  string `json:"rel"`
	HRef string `json:"href"`
}

func HandleNodeInfoWellKnown(w http.ResponseWriter, r *http.Request) {
	nodeInfoLinks := []NodeInfoWellKnownLink{
		{
			Rel:  "http://nodeinfo.diaspora.software/ns/schema/2.0",
			HRef: fmt.Sprintf("https://%s/nodeinfo/2.0.json", r.Host),
		},
	}

	nodeInfo := NodeInfoWellKnown{
		Links: nodeInfoLinks,
	}

	js, err := json.Marshal(nodeInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
