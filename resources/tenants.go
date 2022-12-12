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

// GenerateTenants generates terraform code for Morpheus tenants
func GenerateTenants(client *morpheus.Client) (output []string) {
	log.Println("generating tenants...")
	response, err := client.ListTenants(&morpheus.Request{
		QueryParams: map[string]string{"max": "1000"},
	})
	if err != nil {
		log.Println(err)
	}
	result := response.Result.(*morpheus.ListTenantsResult)
	tenants := result.Accounts
	for _, v := range *tenants {
		// create new empty hcl file object
		hclFile := hclwrite.NewEmptyFile()

		// initialize the body of the new file object
		rootBody := hclFile.Body()
		title := utils.GenerateResourceName(v.Name)
		provider := rootBody.AppendNewBlock("resource",
			[]string{"morpheus_tenant", title})
		providerBody := provider.Body()

		providerBody.SetAttributeValue("name", cty.StringVal(v.Name))
		if v.Description != "" {
			providerBody.SetAttributeValue("description", cty.StringVal(v.Description))
		}
		providerBody.SetAttributeValue("enabled", cty.BoolVal(v.Active))
		if v.Subdomain != "" {
			providerBody.SetAttributeValue("subdomain", cty.StringVal(v.Subdomain))
		}
		//		providerBody.SetAttributeValue("base_role_id", cty.NumberIntVal(v.))
		providerBody.SetAttributeValue("currency", cty.StringVal(v.Currency))
		if v.AccountNumber != "" {
			providerBody.SetAttributeValue("account_number", cty.StringVal(v.AccountNumber))
		}
		if v.AccountName != "" {
			providerBody.SetAttributeValue("account_name", cty.StringVal(v.AccountName))
		}
		if v.CustomerNumber != "" {
			providerBody.SetAttributeValue("customer_number", cty.StringVal(v.CustomerNumber))
		}
		hcloutput := string(hclFile.Bytes())
		output = append(output, hcloutput)
	}
	v := strings.Join(output, "\n")
	err = os.WriteFile("generated/tenants.tf", []byte(v), 0644)
	if err != nil {
		log.Println(err)
	}
	return output
}
