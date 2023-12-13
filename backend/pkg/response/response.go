package response

import (
	"encoding/xml"
	"net/http"

	"github.com/goccy/go-json"
	"sigs.k8s.io/yaml"
)

func Respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	w.WriteHeader(status)

	// Get the Accept header from the request
	accept := r.Header.Get("Accept")

	var resp []byte
	contentType := "application/json"
	// If the Accept header is empty, default to JSON
	if accept == "application/xml" {
		contentType = "application/xml"
		tmp := struct {
			XMLName xml.Name `xml:"response"`
			data    interface{}
		}{data: data}
		resp, _ = xml.Marshal(tmp)
	} else if accept == "text/x-yaml" || accept == "application/x-yaml" || accept == "application/yaml" {
		contentType = accept
		resp, _ = yaml.Marshal(data)
	} else {
		resp, _ = json.Marshal(data)
	}
	w.Header().Set("Content-Type", contentType)
	w.Write([]byte(resp))
}
