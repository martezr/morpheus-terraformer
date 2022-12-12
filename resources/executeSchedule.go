package resources

import (
	"log"
	"os"
	"strings"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// GenerateExecuteSchedules generates terraform code for Morpheus execute schedules
func GenerateExecuteSchedules(client *morpheus.Client) (output []string) {
	log.Println("generating execute schedules...")
	response, err := client.ListExecuteSchedules(&morpheus.Request{
		QueryParams: map[string]string{"max": "1000"},
	})
	if err != nil {
		log.Println(err)
	}
	result := response.Result.(*morpheus.ListExecuteSchedulesResult)
	executeSchedules := result.ExecuteSchedules
	for _, v := range *executeSchedules {
		// create new empty hcl file object
		hclFile := hclwrite.NewEmptyFile()

		// initialize the body of the new file object
		rootBody := hclFile.Body()
		title := strings.ReplaceAll(v.Name, " ", "_")
		title = strings.ToLower(title)
		provider := rootBody.AppendNewBlock("resource",
			[]string{"morpheus_execute_schedule", title})
		providerBody := provider.Body()

		providerBody.SetAttributeValue("name", cty.StringVal(v.Name))
		providerBody.SetAttributeValue("description", cty.StringVal(v.Desription))
		providerBody.SetAttributeValue("enabled", cty.BoolVal(v.Enabled))
		providerBody.SetAttributeValue("time_zone", cty.StringVal(v.TimeZone))
		providerBody.SetAttributeValue("schedule", cty.StringVal(v.Cron))
		hcloutput := string(hclFile.Bytes())
		output = append(output, hcloutput)
	}
	v := strings.Join(output, "\n")
	err = os.WriteFile("generated/executeSchedules.tf", []byte(v), 0644)
	if err != nil {
		log.Println(err)
	}
	return output
}
