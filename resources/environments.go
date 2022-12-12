package resources

import (
	"log"
	"os"
	"strings"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/martezr/morpheus-terraformer/utils"
	"github.com/zclconf/go-cty/cty"
)

// GenerateEnvironments generates terraform code for Morpheus environments
func GenerateEnvironments(client *morpheus.Client) (output []string) {
	log.Println("generating environments...")
	response, err := client.ListEnvironments(&morpheus.Request{
		QueryParams: map[string]string{"max": "1000"},
	})
	if err != nil {
		log.Println(err)
	}
	result := response.Result.(*morpheus.ListEnvironmentsResult)
	environments := result.Environments
	for _, v := range *environments {
		// create new empty hcl file object
		hclFile := hclwrite.NewEmptyFile()

		// initialize the body of the new file object
		rootBody := hclFile.Body()
		title := utils.GenerateResourceName(v.Name)
		provider := rootBody.AppendNewBlock("resource",
			[]string{"morpheus_environment", title})
		providerBody := provider.Body()

		providerBody.SetAttributeValue("name", cty.StringVal(v.Name))
		providerBody.SetAttributeValue("description", cty.StringVal(v.Description))
		providerBody.SetAttributeValue("code", cty.StringVal(v.Code))
		providerBody.SetAttributeValue("active", cty.BoolVal(v.Active))

		hcloutput := string(hclFile.Bytes())
		output = append(output, hcloutput)
	}
	v := strings.Join(output, "\n")
	err = os.WriteFile("generated/environments.tf", []byte(v), 0644)
	if err != nil {
		log.Println(err)
	}
	return output
}
