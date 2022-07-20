package resources

import (
	"log"
	"os"
	"strings"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func GenerateOptionLists(client *morpheus.Client) {
	log.Println("generating option lists...")
	response, err := client.ListOptionLists(&morpheus.Request{
		QueryParams: map[string]string{"max": "500"},
	})
	if err != nil {
		log.Println(err)
	}

	result := response.Result.(*morpheus.ListOptionListsResult)
	optionLists := result.OptionLists
	var manualOptionLists []string
	var restOptionLists []string

	for _, v := range *optionLists {
		switch v.Type {
		case "manual":
			manualOptionLists = append(manualOptionLists, generateManualOptionLists(v))
		case "api":
			//selectOptionTypes = append(selectOptionTypes, generateSelecListOptionTypes(v))
		case "rest":
			restOptionLists = append(restOptionLists, generateRestOptionLists(v))
		}
	}

	// Write Text Option Types
	v := strings.Join(manualOptionLists, "\n")
	err = os.WriteFile("generated/manualOptionLists.tf", []byte(v), 0644)
	if err != nil {
		log.Println(err)
	}

	// Write Text Option Types
	restData := strings.Join(restOptionLists, "\n")
	err = os.WriteFile("generated/restOptionLists.tf", []byte(restData), 0644)
	if err != nil {
		log.Println(err)
	}
}

func generateManualOptionLists(resource morpheus.OptionList) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := strings.ReplaceAll(resource.Name, " ", "_")
	title = strings.ToLower(title)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_manual_option_list", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("description", cty.StringVal(resource.Description))
	data := strings.Split(resource.InitialDataset, "\n")
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

	providerBody.SetAttributeRaw("dataset", out)
	providerBody.SetAttributeValue("real_time", cty.BoolVal(resource.RealTime))
	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateRestOptionLists(resource morpheus.OptionList) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := strings.ReplaceAll(resource.Name, " ", "_")
	title = strings.ToLower(title)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_rest_option_list", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("description", cty.StringVal(resource.Description))
	providerBody.SetAttributeValue("visibility", cty.StringVal(resource.Visibility))
	providerBody.SetAttributeValue("source_url", cty.StringVal(resource.SourceURL))
	providerBody.SetAttributeValue("real_time", cty.BoolVal(resource.RealTime))
	providerBody.SetAttributeValue("ignore_ssl_errors", cty.BoolVal(resource.IgnoreSSLErrors))
	providerBody.SetAttributeValue("initial_dataset", cty.StringVal(resource.InitialDataset))
	providerBody.SetAttributeValue("translation_script", cty.StringVal(resource.TranslationScript))
	providerBody.SetAttributeValue("request_script", cty.StringVal(resource.RequestScript))
	providerBody.SetAttributeValue("source_method", cty.StringVal(resource.SourceMethod))

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}
