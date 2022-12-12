package resources

import (
	"log"
	"os"
	"strings"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/martezr/morpheus-terraformer/utils"
	"github.com/zclconf/go-cty/cty"
)

// GenerateWikis generates terraform code for Morpheus wikis
func GenerateWikis(client *morpheus.Client) (output []string) {
	log.Println("generating wikis...")
	response, err := client.ListWikis(&morpheus.Request{
		QueryParams: map[string]string{"max": "1000"},
	})
	if err != nil {
		log.Println(err)
	}
	result := response.Result.(*morpheus.ListWikisResult)
	wikis := result.Wikis
	for _, v := range *wikis {
		// create new empty hcl file object
		hclFile := hclwrite.NewEmptyFile()

		// initialize the body of the new file object
		rootBody := hclFile.Body()
		title := utils.GenerateResourceName(v.Name)
		provider := rootBody.AppendNewBlock("resource",
			[]string{"morpheus_wiki", title})
		providerBody := provider.Body()
		providerBody.SetAttributeValue("name", cty.StringVal(v.Name))
		providerBody.SetAttributeValue("category", cty.StringVal(v.Category))

		data := strings.Split(v.Content, "\n")
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

		providerBody.SetAttributeRaw("content", out)

		hcloutput := string(hclFile.Bytes())
		output = append(output, hcloutput)
	}
	v := strings.Join(output, "\n")
	err = os.WriteFile("generated/wikis.tf", []byte(v), 0644)
	if err != nil {
		log.Println(err)
	}
	return output
}
