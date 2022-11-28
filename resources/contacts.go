package resources

import (
	"log"
	"os"
	"strings"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/martezr/morpheus-terraformer/utils"
	"github.com/zclconf/go-cty/cty"
)

// GenerateContacts generates terraform code for Morpheus contacts
func GenerateContacts(client *morpheus.Client) (output []string) {
	log.Println("generating contacts...")
	response, err := client.ListContacts(&morpheus.Request{
		QueryParams: map[string]string{"max": "500"},
	})
	if err != nil {
		log.Println(err)
	}
	result := response.Result.(*morpheus.ListContactsResult)
	contacts := result.Contacts
	for _, v := range *contacts {
		// create new empty hcl file object
		hclFile := hclwrite.NewEmptyFile()

		// initialize the body of the new file object
		rootBody := hclFile.Body()
		title := utils.GenerateResourceName(v.Name)
		provider := rootBody.AppendNewBlock("resource",
			[]string{"morpheus_contact", title})
		providerBody := provider.Body()
		providerBody.SetAttributeValue("name", cty.StringVal(v.Name))
		if v.EmailAddress != "" {
			providerBody.SetAttributeValue("email_address", cty.StringVal(v.EmailAddress))
		}
		if v.SmsAddress != "" {
			providerBody.SetAttributeValue("mobile_number", cty.StringVal(v.SmsAddress))
		}

		hcloutput := string(hclFile.Bytes())
		output = append(output, hcloutput)
	}
	v := strings.Join(output, "\n")
	err = os.WriteFile("generated/contacts.tf", []byte(v), 0644)
	if err != nil {
		log.Println(err)
	}
	return output
}
