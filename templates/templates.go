package templates

import (
	"bytes"
	"github.com/diggs/glog"
	template_ "text/template"
)

func Template(data interface{}, template string) (string, error) {

	// TODO
	//  Keep compiled templates in memory
	//  Name templates appropriately (helps with error reporting)
	tmpl, err := template_.New("temp").Parse(template)
	if err != nil {
		return "", err
	}

	output := new(bytes.Buffer)
	err = tmpl.Execute(output, data)
	if err != nil {
		return "", err
	}
	outputString := output.String()

	glog.Debugf("Template result: %s", outputString)

	return outputString, nil
}
