package resources

import (
	"log"
	"os"
	"strings"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func GenerateTenants(client *morpheus.Client) (output []string) {
	log.Println("generating tenants...")
	response, err := client.ListTenants(&morpheus.Request{
		QueryParams: map[string]string{"max": "300"},
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
		title := strings.ReplaceAll(v.Name, " ", "_")
		title = strings.ToLower(title)
		provider := rootBody.AppendNewBlock("resource",
			[]string{"morpheus_tenant", title})
		providerBody := provider.Body()

		providerBody.SetAttributeValue("name", cty.StringVal(v.Name))
		providerBody.SetAttributeValue("description", cty.StringVal(v.Description))
		providerBody.SetAttributeValue("enabled", cty.BoolVal(v.Active))
		providerBody.SetAttributeValue("subdomain", cty.StringVal(v.Subdomain))
		//		providerBody.SetAttributeValue("base_role_id", cty.NumberIntVal(v.))
		providerBody.SetAttributeValue("currency", cty.StringVal(v.Currency))
		providerBody.SetAttributeValue("account_number", cty.StringVal(v.AccountNumber))
		providerBody.SetAttributeValue("account_name", cty.StringVal(v.AccountName))
		providerBody.SetAttributeValue("customer_number", cty.StringVal(v.CustomerNumber))

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
