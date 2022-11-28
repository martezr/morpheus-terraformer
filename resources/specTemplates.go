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

// GenerateSpecTemplates generates terraform code for Morpheus spec templates
func GenerateSpecTemplates(client *morpheus.Client) {
	log.Println("generating spec templates...")
	var generatedSpecTemplates int64
	response, err := client.ListSpecTemplates(&morpheus.Request{
		QueryParams: map[string]string{"max": "1000"},
	})
	if err != nil {
		log.Println(err)
	}

	result := response.Result.(*morpheus.ListSpecTemplatesResult)
	specTemplates := result.SpecTemplates
	baseDir := "generated/spectemplates"
	if _, err := os.Stat(baseDir); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(baseDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	var armSpecTemplates []string
	var cloudFormationSpecTemplates []string
	var helmSpecTemplates []string
	var kubernetesSpecTemplates []string
	var terraformSpecTemplates []string

	for _, v := range *specTemplates {
		// Ignore system created file templates
		if v.Account.ID != 0 {
			generatedSpecTemplates++
			switch v.Type.Code {
			case "arm":
				armSpecTemplates = append(armSpecTemplates, generateArmSpecTemplates(v))
			case "cloudFormation":
				cloudFormationSpecTemplates = append(cloudFormationSpecTemplates, generateCloudFormationSpecTemplates(v))
			case "helm":
				helmSpecTemplates = append(helmSpecTemplates, generateHelmSpecTemplates(v))
			case "kubernetes":
				kubernetesSpecTemplates = append(kubernetesSpecTemplates, generateKubernetesSpecTemplates(v))
			case "terraform":
				terraformSpecTemplates = append(terraformSpecTemplates, generateTerraformSpecTemplates(v))
			}
		}
	}

	if len(armSpecTemplates) > 0 {
		v := strings.Join(armSpecTemplates, "\n")
		err = os.WriteFile("generated/armSpecTemplates.tf", []byte(v), 0644)
		if err != nil {
			log.Println(err)
		}
	}

	if len(cloudFormationSpecTemplates) > 0 {
		selectData := strings.Join(cloudFormationSpecTemplates, "\n")
		err = os.WriteFile("generated/cloudFormationSpecTemplates.tf", []byte(selectData), 0644)
		if err != nil {
			log.Println(err)
		}
	}

	if len(helmSpecTemplates) > 0 {
		hiddenData := strings.Join(helmSpecTemplates, "\n")
		err = os.WriteFile("generated/helmSpecTemplates.tf", []byte(hiddenData), 0644)
		if err != nil {
			log.Println(err)
		}
	}

	if len(kubernetesSpecTemplates) > 0 {
		passwordData := strings.Join(kubernetesSpecTemplates, "\n")
		err = os.WriteFile("generated/kubernetesSpecTemplates.tf", []byte(passwordData), 0644)
		if err != nil {
			log.Println(err)
		}
	}

	if len(terraformSpecTemplates) > 0 {
		numberData := strings.Join(terraformSpecTemplates, "\n")
		err = os.WriteFile("generated/terraformSpecTemplates.tf", []byte(numberData), 0644)
		if err != nil {
			log.Println(err)
		}
	}
	log.Println("processed spec templates:", generatedSpecTemplates)
}

func generateArmSpecTemplates(resource morpheus.SpecTemplate) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_arm_spec_template", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("source_type", cty.StringVal(resource.File.Sourcetype))

	switch resource.File.Sourcetype {
	case "local":
		var out hclwrite.Tokens
		var scriptContent hclwrite.Token
		relativeFilePath := "generated/scripttemplates/" + title + ".json"
		payload := fmt.Sprintf(" file(\"%s\")", relativeFilePath)
		scriptContent.Bytes = []byte(payload)
		out = append(out, &scriptContent)
		providerBody.SetAttributeRaw("spec_content", out)

		filePath, _ := filepath.Abs("generated/spectemplates/" + title + ".json")
		err := os.WriteFile(filePath, []byte(resource.File.Content), 0644)
		if err != nil {
			log.Println(err)
		}

		hcloutput := string(hclFile.Bytes())
		return hcloutput
	case "url":
		providerBody.SetAttributeValue("spec_path", cty.StringVal(resource.File.ContentPath))
	case "git":
		providerBody.SetAttributeValue("spec_path", cty.StringVal(resource.File.ContentPath))
		providerBody.SetAttributeValue("repository_id", cty.NumberIntVal(resource.File.Repository.ID))
		providerBody.SetAttributeValue("version_ref", cty.StringVal(resource.File.ContentRef))
	}

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateCloudFormationSpecTemplates(resource morpheus.SpecTemplate) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_cloud_formation_spec_template", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("source_type", cty.StringVal(resource.File.Sourcetype))

	switch resource.File.Sourcetype {
	case "local":
		providerBody.SetAttributeValue("spec_content", cty.StringVal(resource.File.Content))
		var out hclwrite.Tokens
		var scriptContent hclwrite.Token
		relativeFilePath := "generated/scripttemplates/" + title + ".json"
		payload := fmt.Sprintf(" file(\"%s\")", relativeFilePath)
		scriptContent.Bytes = []byte(payload)
		out = append(out, &scriptContent)
		providerBody.SetAttributeRaw("spec_content", out)

		filePath, _ := filepath.Abs("generated/spectemplates/" + title + ".json")
		err := os.WriteFile(filePath, []byte(resource.File.Content), 0644)
		if err != nil {
			log.Println(err)
		}
	case "url":
		providerBody.SetAttributeValue("spec_path", cty.StringVal(resource.File.ContentPath))
	case "git":
		providerBody.SetAttributeValue("spec_path", cty.StringVal(resource.File.ContentPath))
		providerBody.SetAttributeValue("repository_id", cty.NumberIntVal(resource.File.Repository.ID))
		providerBody.SetAttributeValue("version_ref", cty.StringVal(resource.File.ContentRef))
	}

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateHelmSpecTemplates(resource morpheus.SpecTemplate) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_helm_spec_template", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("source_type", cty.StringVal(resource.File.Sourcetype))

	switch resource.File.Sourcetype {
	case "local":
		providerBody.SetAttributeValue("spec_content", cty.StringVal(resource.File.Content))
	case "url":
		providerBody.SetAttributeValue("spec_path", cty.StringVal(resource.File.ContentPath))
	case "git":
		providerBody.SetAttributeValue("spec_path", cty.StringVal(resource.File.ContentPath))
		providerBody.SetAttributeValue("repository_id", cty.NumberIntVal(resource.File.Repository.ID))
		providerBody.SetAttributeValue("version_ref", cty.StringVal(resource.File.ContentRef))
	}

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateKubernetesSpecTemplates(resource morpheus.SpecTemplate) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_kubernetes_spec_template", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("source_type", cty.StringVal(resource.File.Sourcetype))

	switch resource.File.Sourcetype {
	case "local":
		providerBody.SetAttributeValue("spec_content", cty.StringVal(resource.File.Content))
		var out hclwrite.Tokens
		var scriptContent hclwrite.Token
		relativeFilePath := "generated/scripttemplates/" + title + ".yaml"
		payload := fmt.Sprintf(" file(\"%s\")", relativeFilePath)
		scriptContent.Bytes = []byte(payload)
		out = append(out, &scriptContent)
		providerBody.SetAttributeRaw("spec_content", out)

		filePath, _ := filepath.Abs("generated/spectemplates/" + title + ".yaml")
		err := os.WriteFile(filePath, []byte(resource.File.Content), 0644)
		if err != nil {
			log.Println(err)
		}
	case "url":
		providerBody.SetAttributeValue("spec_path", cty.StringVal(resource.File.ContentPath))
	case "git":
		providerBody.SetAttributeValue("spec_path", cty.StringVal(resource.File.ContentPath))
		providerBody.SetAttributeValue("repository_id", cty.NumberIntVal(resource.File.Repository.ID))
		providerBody.SetAttributeValue("version_ref", cty.StringVal(resource.File.ContentRef))
	}

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateTerraformSpecTemplates(resource morpheus.SpecTemplate) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_terraform_spec_template", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("source_type", cty.StringVal(resource.File.Sourcetype))

	switch resource.File.Sourcetype {
	case "local":
		var out hclwrite.Tokens
		var scriptContent hclwrite.Token
		relativeFilePath := "generated/scripttemplates/" + title + ".tf"
		payload := fmt.Sprintf(" file(\"%s\")", relativeFilePath)
		scriptContent.Bytes = []byte(payload)
		out = append(out, &scriptContent)
		providerBody.SetAttributeRaw("spec_content", out)

		filePath, _ := filepath.Abs("generated/spectemplates/" + title + ".tf")
		err := os.WriteFile(filePath, []byte(resource.File.Content), 0644)
		if err != nil {
			log.Println(err)
		}
	case "url":
		providerBody.SetAttributeValue("spec_path", cty.StringVal(resource.File.ContentPath))
	case "git":
		providerBody.SetAttributeValue("spec_path", cty.StringVal(resource.File.ContentPath))
		providerBody.SetAttributeValue("repository_id", cty.NumberIntVal(resource.File.Repository.ID))
		providerBody.SetAttributeValue("version_ref", cty.StringVal(resource.File.ContentRef))
	}

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}
