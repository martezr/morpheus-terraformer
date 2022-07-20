package resources

import (
	"log"
	"os"
	"strings"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// GenerateGroups generates terraform code for Morpheus groups
func GenerateGroups(client *morpheus.Client) (output []string) {
	log.Println("generating groups...")
	response, err := client.ListGroups(&morpheus.Request{
		QueryParams: map[string]string{"max": "300"},
	})
	if err != nil {
		log.Println(err)
	}
	result := response.Result.(*morpheus.ListGroupsResult)
	groups := result.Groups
	for _, v := range *groups {
		// create new empty hcl file object
		hclFile := hclwrite.NewEmptyFile()

		// initialize the body of the new file object
		rootBody := hclFile.Body()
		title := strings.ReplaceAll(v.Name, " ", "_")
		title = strings.ToLower(title)
		provider := rootBody.AppendNewBlock("resource",
			[]string{"morpheus_group", title})
		providerBody := provider.Body()

		providerBody.SetAttributeValue("name", cty.StringVal(v.Name))
		providerBody.SetAttributeValue("code", cty.StringVal(v.Code))
		providerBody.SetAttributeValue("location", cty.StringVal(v.Location))

		hcloutput := string(hclFile.Bytes())
		output = append(output, hcloutput)
	}
	v := strings.Join(output, "\n")
	err = os.WriteFile("generated/groups.tf", []byte(v), 0644)
	if err != nil {
		log.Println(err)
	}
	return output
}
