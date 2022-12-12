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

// GenerateFileTemplates generates terraform code for Morpheus file templates
func GenerateFileTemplates(client *morpheus.Client) (output []string) {
	log.Println("generating file templates...")
	var generatedFileTemplates int64
	response, err := client.ListFileTemplates(&morpheus.Request{
		QueryParams: map[string]string{"max": "1000"},
	})
	if err != nil {
		log.Println(err)
	}
	result := response.Result.(*morpheus.ListFileTemplatesResult)
	fileTemplates := result.FileTemplates
	baseDir := "generated/filetemplates"
	if _, err := os.Stat(baseDir); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(baseDir, 0755)
		if err != nil {
			panic(err)
		}
	}
	for _, v := range *fileTemplates {
		// Ignore system created file templates
		if v.Account.ID != 0 {
			generatedFileTemplates++
			// create new empty hcl file object
			hclFile := hclwrite.NewEmptyFile()

			// initialize the body of the new file object
			rootBody := hclFile.Body()
			title := utils.GenerateResourceName(v.Name)
			provider := rootBody.AppendNewBlock("resource",
				[]string{"morpheus_file_template", title})
			providerBody := provider.Body()
			providerBody.SetAttributeValue("name", cty.StringVal(v.Name))
			providerBody.SetAttributeValue("file_name", cty.StringVal(v.FileName))
			providerBody.SetAttributeValue("file_path", cty.StringVal(v.FilePath))
			providerBody.SetAttributeValue("phase", cty.StringVal(v.TemplatePhase))
			if v.FileOwner != "" {
				providerBody.SetAttributeValue("file_owner", cty.StringVal(v.FileOwner))
			}
			if v.SettingName != "" {
				providerBody.SetAttributeValue("setting_name", cty.StringVal(v.SettingName))
			}
			if v.SettingCategory != "" {
				providerBody.SetAttributeValue("setting_category", cty.StringVal(v.SettingCategory))
			}

			var out hclwrite.Tokens
			var scriptContent hclwrite.Token
			relativeFilePath := "generated/filetemplates/" + v.FileName
			payload := fmt.Sprintf(" file(\"%s\")", relativeFilePath)
			scriptContent.Bytes = []byte(payload)
			out = append(out, &scriptContent)
			providerBody.SetAttributeRaw("file_content", out)

			hcloutput := string(hclFile.Bytes())
			output = append(output, hcloutput)
			filePath, _ := filepath.Abs("generated/filetemplates/" + v.FileName)
			err = os.WriteFile(filePath, []byte(v.Template), 0644)
			if err != nil {
				log.Println(err)
			}

		}
	}
	v := strings.Join(output, "\n")
	err = os.WriteFile("generated/fileTemplates.tf", []byte(v), 0644)
	if err != nil {
		log.Println(err)
	}
	log.Println("processed file templates:", generatedFileTemplates)
	return output
}
