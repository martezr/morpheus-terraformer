package resources

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// GenerateWorkflows generates terraform code for Morpheus workflows
func GenerateWorkflows(client *morpheus.Client) {
	response, err := client.ListTaskSets(&morpheus.Request{
		QueryParams: map[string]string{"max": "500"},
	})
	if err != nil {
		log.Println(err)
	}
	var workflowResponse WorkflowsPayload
	json.Unmarshal(response.Body, &workflowResponse)
	var operationalWorkflows []string

	for _, v := range workflowResponse.TaskSets {
		switch v.Type {
		case "operation":
			operationalWorkflows = append(operationalWorkflows, generateOperationalWorkflows(v))
			generateOperationalWorkflows(v)
		case "provision":
			//generateProvisioningWorkflows(v)
		}
	}

	operationalWorkflowData := strings.Join(operationalWorkflows, "\n")
	err = os.WriteFile("generated/operationalWorkflows.tf", []byte(operationalWorkflowData), 0644)
	if err != nil {
		log.Println(err)
	}
}

func generateOperationalWorkflows(resource Workflow) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := strings.ReplaceAll(resource.Name, " ", "_")
	title = strings.ToLower(title)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_operational_workflow", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("description", cty.StringVal(resource.Description))
	providerBody.SetAttributeValue("visibility", cty.StringVal(resource.Visibility))
	if resource.Platform == "" {
		providerBody.SetAttributeValue("platform", cty.StringVal("all"))
	} else {
		providerBody.SetAttributeValue("platform", cty.StringVal(resource.Platform))
	}
	providerBody.SetAttributeValue("allow_custom_config", cty.BoolVal(resource.Allowcustomconfig))
	var tasks []cty.Value
	for _, v := range resource.Tasksettasks {
		taskRef := strings.ReplaceAll(v.Task.Name, " ", "_")
		taskRef = strings.ToLower(taskRef)
		taskType := v.Task.Tasktype.Code
		var resourceOut string
		switch taskType {
		case "jythonTask":
			resourceOut = "morpheus_python_script_task" + "." + taskRef + ".id"
			tasks = append(tasks, cty.StringVal(resourceOut))
		}
	}
	if len(tasks) == 0 {
		taskid := cty.String
		providerBody.SetAttributeValue("task_ids", cty.ListValEmpty(taskid))
	} else {
		providerBody.SetAttributeValue("task_ids", cty.ListVal(tasks))
	}
	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

// Workflow defines a Morpheus workflow object
type Workflow struct {
	Accountid         int       `json:"accountId"`
	Allowcustomconfig bool      `json:"allowCustomConfig"`
	Datecreated       time.Time `json:"dateCreated"`
	Description       string    `json:"description"`
	ID                int       `json:"id"`
	Lastupdated       time.Time `json:"lastUpdated"`
	Name              string    `json:"name"`
	Optiontypes       []struct {
		Advanced bool   `json:"advanced"`
		Code     string `json:"code"`
		Config   struct {
		} `json:"config,omitempty"`
		Contextualdefault     bool        `json:"contextualDefault"`
		Creatable             bool        `json:"creatable"`
		Defaultvalue          interface{} `json:"defaultValue"`
		Dependsoncode         interface{} `json:"dependsOnCode"`
		Description           string      `json:"description"`
		Displayorder          int         `json:"displayOrder"`
		Displayvalueondetails bool        `json:"displayValueOnDetails"`
		Editable              bool        `json:"editable"`
		Enabled               bool        `json:"enabled"`
		Exportmeta            bool        `json:"exportMeta"`
		Fieldaddon            interface{} `json:"fieldAddOn"`
		Fieldclass            interface{} `json:"fieldClass"`
		Fieldcode             interface{} `json:"fieldCode"`
		Fieldcomponent        interface{} `json:"fieldComponent"`
		Fieldcontext          string      `json:"fieldContext"`
		Fieldgroup            interface{} `json:"fieldGroup"`
		Fieldinput            interface{} `json:"fieldInput"`
		Fieldlabel            string      `json:"fieldLabel"`
		Fieldname             string      `json:"fieldName"`
		Helpblock             interface{} `json:"helpBlock"`
		Helpblockfieldcode    interface{} `json:"helpBlockFieldCode"`
		ID                    int         `json:"id"`
		Name                  string      `json:"name"`
		Noblank               bool        `json:"noBlank"`
		Optionlist            interface{} `json:"optionList"`
		Optionsource          interface{} `json:"optionSource"`
		Optionsourcetype      interface{} `json:"optionSourceType"`
		Placeholder           interface{} `json:"placeHolder"`
		Requireoncode         interface{} `json:"requireOnCode"`
		Required              bool        `json:"required"`
		Showoncreate          bool        `json:"showOnCreate"`
		Showonedit            bool        `json:"showOnEdit"`
		Type                  string      `json:"type"`
		Verifypattern         interface{} `json:"verifyPattern"`
		Visibleoncode         interface{} `json:"visibleOnCode"`
		Wrapperclass          interface{} `json:"wrapperClass"`
	} `json:"optionTypes"`
	Platform     string `json:"platform"`
	Tasksettasks []struct {
		ID   int `json:"id"`
		Task struct {
			Accountid         int       `json:"accountId"`
			Allowcustomconfig bool      `json:"allowCustomConfig"`
			Code              string    `json:"code"`
			Datecreated       time.Time `json:"dateCreated"`
			Executetarget     string    `json:"executeTarget"`
			File              struct {
				Content     string      `json:"content"`
				Contentpath interface{} `json:"contentPath"`
				Contentref  interface{} `json:"contentRef"`
				ID          int         `json:"id"`
				Repository  interface{} `json:"repository"`
				Sourcetype  string      `json:"sourceType"`
			} `json:"file"`
			ID                int         `json:"id"`
			Lastupdated       time.Time   `json:"lastUpdated"`
			Name              string      `json:"name"`
			Resulttype        interface{} `json:"resultType"`
			Retrycount        int         `json:"retryCount"`
			Retrydelayseconds int         `json:"retryDelaySeconds"`
			Retryable         bool        `json:"retryable"`
			Taskoptions       struct {
				Host                     interface{} `json:"host"`
				Localscriptgitid         interface{} `json:"localScriptGitId"`
				Localscriptgitref        interface{} `json:"localScriptGitRef"`
				Password                 interface{} `json:"password"`
				Port                     interface{} `json:"port"`
				Pythonadditionalpackages string      `json:"pythonAdditionalPackages"`
				Pythonargs               interface{} `json:"pythonArgs"`
				Pythonbinary             string      `json:"pythonBinary"`
				Sshkey                   interface{} `json:"sshKey"`
				Username                 interface{} `json:"username"`
			} `json:"taskOptions"`
			Tasktype struct {
				Code string `json:"code"`
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"taskType"`
		} `json:"task"`
		Taskorder int    `json:"taskOrder"`
		Taskphase string `json:"taskPhase"`
	} `json:"taskSetTasks"`
	Tasks      []int  `json:"tasks"`
	Type       string `json:"type"`
	Visibility string `json:"visibility"`
}

// WorkflowsPayload defines a Morpheus workflow objects
type WorkflowsPayload struct {
	TaskSets []Workflow `json:"taskSets"`
	Meta     struct {
		Offset int `json:"offset"`
		Max    int `json:"max"`
		Size   int `json:"size"`
		Total  int `json:"total"`
	} `json:"meta"`
}
