package resources

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// GenerateTasks generates terraform code for Morpheus tasks
func GenerateTasks(client *morpheus.Client) {
	log.Println("generating tasks...")
	response, err := client.ListTasks(&morpheus.Request{
		QueryParams: map[string]string{"max": "500"},
	})
	if err != nil {
		log.Println(err)
	}

	var taskResponse TasksPayload
	json.Unmarshal(response.Body, &taskResponse)
	var ansibleTasks []string
	var pythonTasks []string
	var scriptTasks []string

	for _, v := range taskResponse.Tasks {
		switch v.Tasktype.Code {
		case "ansibleTask":
			ansibleTasks = append(ansibleTasks, generateAnsibleTask(v))
		case "groovyTask":
		case "script":
			scriptTasks = append(scriptTasks, generateShellScriptTask(v))
		case "jythonTask":
			pythonTasks = append(pythonTasks, generatePythonTask(v))
		case "containerScript":
		case "winrmTask":
		case "httpTask":
		case "puppetTask":
		case "email":
		case "containerTemplate":
		case "chefTask":
		}
	}

	if len(ansibleTasks) > 0 {
		ansibleTaskData := strings.Join(ansibleTasks, "\n")
		err = os.WriteFile("generated/ansibleTasks.tf", []byte(ansibleTaskData), 0644)
		if err != nil {
			log.Println(err)
		}
	}

	if len(pythonTasks) > 0 {
		pythonTaskData := strings.Join(pythonTasks, "\n")
		err = os.WriteFile("generated/pythonTasks.tf", []byte(pythonTaskData), 0644)
		if err != nil {
			log.Println(err)
		}
	}

	if len(scriptTasks) > 0 {
		scriptTaskData := strings.Join(scriptTasks, "\n")
		err = os.WriteFile("generated/scriptTasks.tf", []byte(scriptTaskData), 0644)
		if err != nil {
			log.Println(err)
		}
	}
}

func generateAnsibleTask(resource Task) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := strings.ReplaceAll(resource.Name, " ", "_")
	title = strings.ToLower(title)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_ansible_playbook_task", title})
	providerBody := provider.Body()
	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("code", cty.StringVal(resource.Code))
	providerBody.SetAttributeValue("ansible_repo_id", cty.StringVal(resource.Taskoptions.Ansiblegitid))
	providerBody.SetAttributeValue("git_ref", cty.StringVal(resource.Taskoptions.Ansiblegitref))
	providerBody.SetAttributeValue("playbook", cty.StringVal(resource.Taskoptions.Ansibleplaybook))
	providerBody.SetAttributeValue("tags", cty.StringVal(resource.Taskoptions.Ansibletags))
	providerBody.SetAttributeValue("skip_tags", cty.StringVal(resource.Taskoptions.Ansibleskiptags))
	providerBody.SetAttributeValue("command_options", cty.StringVal(resource.Taskoptions.Ansibleoptions))
	providerBody.SetAttributeValue("execute_target", cty.StringVal(resource.Executetarget))
	providerBody.SetAttributeValue("retryable", cty.BoolVal(resource.Retryable))

	//	providerBody.SetAttributeValue("retry_delay_seconds", cty.StringVal(resource.Executetarget))
	providerBody.SetAttributeValue("allow_custom_config", cty.BoolVal(resource.Allowcustomconfig))

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generatePythonTask(resource Task) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := strings.ReplaceAll(resource.Name, " ", "_")
	title = strings.ToLower(title)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_python_script_task", title})
	providerBody := provider.Body()
	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("code", cty.StringVal(resource.Code))
	providerBody.SetAttributeValue("source_type", cty.StringVal(resource.File.Sourcetype))
	data := strings.Split(resource.File.Content, "\n")
	var out hclwrite.Tokens
	var eofStart hclwrite.Token
	eofStart.Type = hclsyntax.TokenOHeredoc
	eofStart.Bytes = []byte("<<TFEOF\n")
	out = append(out, &eofStart)

	for _, input := range data {
		var token hclwrite.Token
		token.Type = hclsyntax.TokenStringLit
		token.Bytes = []byte("\t" + input + "\n")
		out = append(out, &token)
	}

	var eofEnd hclwrite.Token
	eofEnd.Type = hclsyntax.TokenCHeredoc
	eofEnd.Bytes = []byte("\tTFEOF")
	out = append(out, &eofEnd)

	providerBody.SetAttributeRaw("script_content", out)
	providerBody.SetAttributeValue("command_arguments", cty.StringVal(resource.Taskoptions.Pythonargs))
	providerBody.SetAttributeValue("retryable", cty.BoolVal(resource.Retryable))
	providerBody.SetAttributeValue("python_binary", cty.StringVal(resource.Taskoptions.Pythonbinary))
	providerBody.SetAttributeValue("additional_packages", cty.StringVal(resource.Taskoptions.Pythonadditionalpackages))
	providerBody.SetAttributeValue("retry_count", cty.NumberIntVal(int64(resource.Retrycount)))
	providerBody.SetAttributeValue("retry_delay_seconds", cty.NumberIntVal(int64(resource.Retrydelayseconds)))

	//	providerBody.SetAttributeValue("retry_delay_seconds", cty.StringVal(resource.Executetarget))
	providerBody.SetAttributeValue("allow_custom_config", cty.BoolVal(resource.Allowcustomconfig))

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateShellScriptTask(resource Task) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := strings.ReplaceAll(resource.Name, " ", "_")
	title = strings.ToLower(title)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_shell_script_task", title})
	providerBody := provider.Body()
	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("code", cty.StringVal(resource.Code))
	providerBody.SetAttributeValue("source_type", cty.StringVal(resource.File.Sourcetype))

	data := strings.Split(resource.File.Content, "\n")
	var out hclwrite.Tokens
	var eofStart hclwrite.Token
	eofStart.Type = hclsyntax.TokenOHeredoc
	eofStart.Bytes = []byte("<<TFEOF\n")
	out = append(out, &eofStart)

	for _, input := range data {
		var token hclwrite.Token
		token.Type = hclsyntax.TokenStringLit
		token.Bytes = []byte("\t" + input + "\n")
		out = append(out, &token)
	}

	var eofEnd hclwrite.Token
	eofEnd.Type = hclsyntax.TokenCHeredoc
	eofEnd.Bytes = []byte("\tTFEOF")
	out = append(out, &eofEnd)

	providerBody.SetAttributeRaw("script_content", out)

	//providerBody.SetAttributeValue("script_content", cty.StringVal(resource.File.Content))
	providerBody.SetAttributeValue("execute_target", cty.StringVal(resource.Executetarget))
	providerBody.SetAttributeValue("retryable", cty.BoolVal(resource.Retryable))
	//	providerBody.SetAttributeValue("retry_delay_seconds", cty.StringVal(resource.Executetarget))
	providerBody.SetAttributeValue("allow_custom_config", cty.BoolVal(resource.Allowcustomconfig))

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

// Task is a Morpheus task object
type Task struct {
	ID        int    `json:"id"`
	Accountid int    `json:"accountId"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	Tasktype  struct {
		ID   int    `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"taskType"`
	File struct {
		ID          int    `json:"id"`
		Sourcetype  string `json:"sourceType"`
		Contentref  string `json:"contentRef"`
		Contentpath string `json:"contentPath"`
		Repository  struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"repository"`
		Content string `json:"content"`
	} `json:"file"`
	Taskoptions struct {
		Ansibleoptions           string      `json:"ansibleOptions"`
		Ansibletags              string      `json:"ansibleTags"`
		Ansibleplaybook          string      `json:"ansiblePlaybook"`
		Ansiblegitref            string      `json:"ansibleGitRef"`
		Ansibleskiptags          string      `json:"ansibleSkipTags"`
		Ansiblegitid             string      `json:"ansibleGitId"`
		Pythonbinary             string      `json:"pythonBinary"`
		Pythonargs               string      `json:"pythonArgs"`
		Pythonadditionalpackages string      `json:"pythonAdditionalPackages"`
		Pythonscript             interface{} `json:"pythonScript"`
	} `json:"taskOptions"`
	Executetarget     string `json:"executeTarget"`
	Retryable         bool   `json:"retryable"`
	Retrycount        int    `json:"retryCount"`
	Retrydelayseconds int    `json:"retryDelaySeconds"`
	Allowcustomconfig bool   `json:"allowCustomConfig"`
}

// TasksPayload is the payload returned for Morpheus task objects
type TasksPayload struct {
	Tasks []Task `json:"tasks"`
	Meta  struct {
		Offset int `json:"offset"`
		Max    int `json:"max"`
		Size   int `json:"size"`
		Total  int `json:"total"`
	} `json:"meta"`
}
