package resources

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/martezr/morpheus-terraformer/utils"
	"github.com/zclconf/go-cty/cty"
)

// GenerateScriptTemplates generates terraform code for Morpheus script templates
func GenerateScriptTemplates(client *morpheus.Client) (output []string) {
	log.Println("generating script templates...")
	var generatedScriptTemplates int64
	response, err := client.ListScriptTemplates(&morpheus.Request{
		QueryParams: map[string]string{"max": "500"},
	})
	if err != nil {
		log.Println(err)
	}
	result := response.Result.(*morpheus.ListScriptTemplatesResult)
	scriptTemplates := result.ScriptTemplates
	baseDir := "generated/scripttemplates"
	if _, err := os.Stat(baseDir); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(baseDir, 0755)
		if err != nil {
			panic(err)
		}
	}
	for _, v := range *scriptTemplates {
		// Ignore system created script templates
		if v.Account.ID != 0 {
			generatedScriptTemplates++
			// create new empty hcl file object
			hclFile := hclwrite.NewEmptyFile()

			// initialize the body of the new file object
			rootBody := hclFile.Body()
			title := utils.GenerateResourceName(v.Name)
			provider := rootBody.AppendNewBlock("resource",
				[]string{"morpheus_script_template", title})
			providerBody := provider.Body()
			providerBody.SetAttributeValue("name", cty.StringVal(v.Name))
			providerBody.SetAttributeValue("script_type", cty.StringVal(v.ScriptType))
			providerBody.SetAttributeValue("script_phase", cty.StringVal(v.ScriptPhase))
			if v.RunAsUser != "" {
				providerBody.SetAttributeValue("run_as_user", cty.StringVal(v.RunAsUser))
			}
			providerBody.SetAttributeValue("sudo", cty.BoolVal(v.SudoUser))
			var scriptExtension string
			if v.ScriptType == "powershell" {
				scriptExtension = ".ps1"
			} else {
				scriptExtension = ".sh"
			}
			var out hclwrite.Tokens
			var scriptContent hclwrite.Token
			relativeFilePath := "generated/scripttemplates/" + title + scriptExtension
			payload := fmt.Sprintf(" file(\"%s\")", relativeFilePath)
			scriptContent.Bytes = []byte(payload)
			out = append(out, &scriptContent)
			providerBody.SetAttributeRaw("script_content", out)

			hcloutput := string(hclFile.Bytes())
			output = append(output, hcloutput)
			filePath, _ := filepath.Abs("generated/scripttemplates/" + title + scriptExtension)
			err = os.WriteFile(filePath, []byte(v.Script), 0644)
			if err != nil {
				log.Println(err)
			}
		}
	}
	v := strings.Join(output, "\n")
	err = os.WriteFile("generated/scriptTemplates.tf", []byte(v), 0644)
	if err != nil {
		log.Println(err)
	}
	log.Println("processed script templates:", generatedScriptTemplates)
	return output
}
