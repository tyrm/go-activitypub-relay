package web

import (
	"html/template"
	"net/http"

	"github.com/tyrm/go-activitypub-relay/models"
)

type IndexTemplate struct {
	Description string
	InstanceList []string
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	if isAccepteType(r, "text/html") {
		tHTML, err := templates.FindString("index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		t, err := template.New("index").Parse(tHTML)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}


		instanceList, err := models.GetApprovedInstances()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		templateVars := IndexTemplate{}

		for _, i := range *instanceList {
			templateVars.InstanceList = append(templateVars.InstanceList, i.Hostname)
		}

		t.Execute(w, &templateVars)
		return
	}


	//w.Header().Set("Content-Type", "application/json")
	//w.Write()
}

