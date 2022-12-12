package resources

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/martezr/morpheus-terraformer/utils"
	"github.com/zclconf/go-cty/cty"
)

// GenerateTasks generates terraform code for Morpheus tasks
func GenerateTasks(client *morpheus.Client) {
	log.Println("generating tasks...")
	response, err := client.ListTasks(&morpheus.Request{
		QueryParams: map[string]string{"max": "1000"},
	})
	if err != nil {
		log.Println(err)
	}

	result := response.Result.(*morpheus.ListTasksResult)
	tasks := result.Tasks

	baseDir := "generated/tasks"
	if _, err := os.Stat(baseDir); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(baseDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	var ansibleTasks []string
	var pythonTasks []string
	var restartTasks []string
	var rubyTasks []string
	var scriptTasks []string
	var writeAttributesTasks []string

	for _, v := range *tasks {
		switch v.TaskType.Code {
		case "ansibleTask":
			ansibleTasks = append(ansibleTasks, generateAnsibleTask(v))
		case "groovyTask":
		case "script":
			scriptTasks = append(scriptTasks, generateShellScriptTask(v))
		case "jythonTask":
			pythonTasks = append(pythonTasks, generatePythonTask(v))
		case "jrubyTask":
			rubyTasks = append(rubyTasks, generateRubyTask(v))
		case "containerScript":
		case "winrmTask":
		case "httpTask":
		case "restartTask":
			restartTasks = append(restartTasks, generateRestartTask(v))
		case "puppetTask":
		case "email":
		case "containerTemplate":
		case "chefTask":
		case "writeAttributes":
			writeAttributesTasks = append(writeAttributesTasks, generateWriteAttributeTask(v))
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

	if len(restartTasks) > 0 {
		restartTaskData := strings.Join(restartTasks, "\n")
		err = os.WriteFile("generated/restartTasks.tf", []byte(restartTaskData), 0644)
		if err != nil {
			log.Println(err)
		}
	}

	if len(rubyTasks) > 0 {
		rubyTaskData := strings.Join(rubyTasks, "\n")
		err = os.WriteFile("generated/rubyTasks.tf", []byte(rubyTaskData), 0644)
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

	if len(writeAttributesTasks) > 0 {
		writeAttributesTaskData := strings.Join(writeAttributesTasks, "\n")
		err = os.WriteFile("generated/writeAttributesTasks.tf", []byte(writeAttributesTaskData), 0644)
		if err != nil {
			log.Println(err)
		}
	}
}

func generateAnsibleTask(resource morpheus.Task) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_ansible_playbook_task", title})
	providerBody := provider.Body()
	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("code", cty.StringVal(resource.Code))
	providerBody.SetAttributeValue("ansible_repo_id", cty.StringVal(resource.TaskOptions.AnsibleGitId))
	providerBody.SetAttributeValue("git_ref", cty.StringVal(resource.TaskOptions.AnsibleGitRef))
	providerBody.SetAttributeValue("playbook", cty.StringVal(resource.TaskOptions.AnsiblePlaybook))
	providerBody.SetAttributeValue("tags", cty.StringVal(resource.TaskOptions.AnsibleTags))
	providerBody.SetAttributeValue("skip_tags", cty.StringVal(resource.TaskOptions.AnsibleSkipTags))
	providerBody.SetAttributeValue("command_options", cty.StringVal(resource.TaskOptions.AnsibleOptions))
	providerBody.SetAttributeValue("execute_target", cty.StringVal(resource.ExecuteTarget))
	providerBody.SetAttributeValue("retryable", cty.BoolVal(resource.Retryable))
	providerBody.SetAttributeValue("retry_delay_seconds", cty.NumberIntVal(resource.RetryDelaySeconds))
	providerBody.SetAttributeValue("allow_custom_config", cty.BoolVal(resource.AllowCustomConfig))

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generatePythonTask(resource morpheus.Task) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_python_script_task", title})
	providerBody := provider.Body()
	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("code", cty.StringVal(resource.Code))
	providerBody.SetAttributeValue("source_type", cty.StringVal(resource.File.SourceType))
	if resource.ResultType != "" {
		providerBody.SetAttributeValue("result_type", cty.StringVal(resource.ResultType))
	}
	switch resource.File.SourceType {
	case "local":
		var out hclwrite.Tokens
		var scriptContent hclwrite.Token
		relativeFilePath := "generated/tasks/" + title + ".py"
		payload := fmt.Sprintf(" file(\"%s\")", relativeFilePath)
		scriptContent.Bytes = []byte(payload)
		out = append(out, &scriptContent)
		providerBody.SetAttributeRaw("script_content", out)

		filePath, _ := filepath.Abs("generated/tasks/" + title + ".py")
		err := os.WriteFile(filePath, []byte(resource.File.Content), 0644)
		if err != nil {
			log.Println(err)
		}
	case "url":
		providerBody.SetAttributeValue("spec_path", cty.StringVal(resource.File.ContentPath))
	case "repository":
		providerBody.SetAttributeValue("spec_path", cty.StringVal(resource.File.ContentPath))
		providerBody.SetAttributeValue("repository_id", cty.NumberIntVal(int64(resource.File.Repository.ID)))
		providerBody.SetAttributeValue("version_ref", cty.StringVal(resource.File.ContentRef))
	}

	providerBody.SetAttributeValue("command_arguments", cty.StringVal(resource.TaskOptions.PythonArgs))
	providerBody.SetAttributeValue("retryable", cty.BoolVal(resource.Retryable))
	providerBody.SetAttributeValue("python_binary", cty.StringVal(resource.TaskOptions.PythonBinary))
	providerBody.SetAttributeValue("additional_packages", cty.StringVal(resource.TaskOptions.PythonAdditionalPackages))
	providerBody.SetAttributeValue("retry_count", cty.NumberIntVal(int64(resource.RetryCount)))
	providerBody.SetAttributeValue("retry_delay_seconds", cty.NumberIntVal(int64(resource.RetryDelaySeconds)))
	providerBody.SetAttributeValue("allow_custom_config", cty.BoolVal(resource.AllowCustomConfig))

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateRestartTask(resource morpheus.Task) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_restart_task", title})
	providerBody := provider.Body()
	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("code", cty.StringVal(resource.Code))
	providerBody.SetAttributeValue("retryable", cty.BoolVal(resource.Retryable))
	providerBody.SetAttributeValue("retry_count", cty.NumberIntVal(int64(resource.RetryCount)))
	providerBody.SetAttributeValue("retry_delay_seconds", cty.NumberIntVal(int64(resource.RetryDelaySeconds)))
	providerBody.SetAttributeValue("allow_custom_config", cty.BoolVal(resource.AllowCustomConfig))

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateRubyTask(resource morpheus.Task) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_ruby_script_task", title})
	providerBody := provider.Body()
	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("code", cty.StringVal(resource.Code))
	providerBody.SetAttributeValue("source_type", cty.StringVal(resource.File.SourceType))
	providerBody.SetAttributeValue("result_type", cty.StringVal(resource.ResultType))
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
	providerBody.SetAttributeValue("retryable", cty.BoolVal(resource.Retryable))
	providerBody.SetAttributeValue("retry_count", cty.NumberIntVal(int64(resource.RetryCount)))
	providerBody.SetAttributeValue("retry_delay_seconds", cty.NumberIntVal(int64(resource.RetryDelaySeconds)))
	providerBody.SetAttributeValue("allow_custom_config", cty.BoolVal(resource.AllowCustomConfig))

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateShellScriptTask(resource morpheus.Task) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_shell_script_task", title})
	providerBody := provider.Body()
	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("code", cty.StringVal(resource.Code))
	providerBody.SetAttributeValue("source_type", cty.StringVal(resource.File.SourceType))
	if resource.ResultType != "" {
		providerBody.SetAttributeValue("result_type", cty.StringVal(resource.ResultType))
	}

	switch resource.File.SourceType {
	case "local":
		var out hclwrite.Tokens
		var scriptContent hclwrite.Token
		relativeFilePath := "generated/tasks/" + title + ".sh"
		payload := fmt.Sprintf(" file(\"%s\")", relativeFilePath)
		scriptContent.Bytes = []byte(payload)
		out = append(out, &scriptContent)
		providerBody.SetAttributeRaw("script_content", out)

		filePath, _ := filepath.Abs("generated/tasks/" + title + ".sh")
		err := os.WriteFile(filePath, []byte(resource.File.Content), 0644)
		if err != nil {
			log.Println(err)
		}
	case "url":
		providerBody.SetAttributeValue("spec_path", cty.StringVal(resource.File.ContentPath))
	case "git":
		providerBody.SetAttributeValue("spec_path", cty.StringVal(resource.File.ContentPath))
		providerBody.SetAttributeValue("repository_id", cty.NumberIntVal(int64(resource.File.Repository.ID)))
		providerBody.SetAttributeValue("version_ref", cty.StringVal(resource.File.ContentRef))
	}

	if resource.TaskOptions.ShellSudo == "on" {
		providerBody.SetAttributeValue("sudo", cty.BoolVal(true))
	} else {
		providerBody.SetAttributeValue("sudo", cty.BoolVal(false))
	}

	providerBody.SetAttributeValue("execute_target", cty.StringVal(resource.ExecuteTarget))
	providerBody.SetAttributeValue("retryable", cty.BoolVal(resource.Retryable))
	providerBody.SetAttributeValue("retry_count", cty.NumberIntVal(int64(resource.RetryCount)))
	providerBody.SetAttributeValue("retry_delay_seconds", cty.NumberIntVal(int64(resource.RetryDelaySeconds)))
	providerBody.SetAttributeValue("allow_custom_config", cty.BoolVal(resource.AllowCustomConfig))

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateWriteAttributeTask(resource morpheus.Task) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_write_attributes_task", title})
	providerBody := provider.Body()
	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("code", cty.StringVal(resource.Code))

	data := strings.Split(resource.TaskOptions.WriteAttributesAttributes, "\n")
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

	providerBody.SetAttributeRaw("attributes", out)
	providerBody.SetAttributeValue("retryable", cty.BoolVal(resource.Retryable))
	providerBody.SetAttributeValue("retry_count", cty.NumberIntVal(int64(resource.RetryCount)))
	providerBody.SetAttributeValue("retry_delay_seconds", cty.NumberIntVal(int64(resource.RetryDelaySeconds)))
	providerBody.SetAttributeValue("allow_custom_config", cty.BoolVal(resource.AllowCustomConfig))

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}
