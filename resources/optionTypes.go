package resources

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/martezr/morpheus-terraformer/utils"
	"github.com/zclconf/go-cty/cty"
)

// GenerateOptionTypes generates terraform code for Morpheus option types
func GenerateOptionTypes(client *morpheus.Client) {
	log.Println("generating option types...")
	response, err := client.ListOptionTypes(&morpheus.Request{
		QueryParams: map[string]string{"max": "500"},
	})
	if err != nil {
		log.Println(err)
	}

	result := response.Result.(*morpheus.ListOptionTypesResult)
	optionTypes := result.OptionTypes
	var textOptionTypes []string
	var selectOptionTypes []string
	var hiddenOptionTypes []string
	var passwordOptionTypes []string
	var numberOptionTypes []string
	var typeAheadOptionTypes []string
	var checkboxOptionTypes []string

	for _, v := range *optionTypes {
		switch v.Type {
		case "text":
			textOptionTypes = append(textOptionTypes, generateTextOptionTypes(v))
		case "select":
			selectOptionTypes = append(selectOptionTypes, generateSelecListOptionTypes(v))
		case "hidden":
			hiddenOptionTypes = append(hiddenOptionTypes, generateHiddenOptionTypes(v))
		case "checkbox":
			checkboxOptionTypes = append(checkboxOptionTypes, generateCheckboxOptionTypes(v))
		case "password":
			passwordOptionTypes = append(passwordOptionTypes, generatePasswordOptionTypes(v))
		case "number":
			numberOptionTypes = append(numberOptionTypes, generateNumberOptionTypes(v))
		case "typeahead":
			typeAheadOptionTypes = append(typeAheadOptionTypes, generateTypeAheadOptionTypes(v))
		}
	}

	if len(textOptionTypes) > 0 {
		v := strings.Join(textOptionTypes, "\n")
		err = os.WriteFile("generated/textOptionTypes.tf", []byte(v), 0644)
		if err != nil {
			log.Println(err)
		}
	}

	if len(selectOptionTypes) > 0 {
		selectData := strings.Join(selectOptionTypes, "\n")
		err = os.WriteFile("generated/selectOptionTypes.tf", []byte(selectData), 0644)
		if err != nil {
			log.Println(err)
		}
	}

	if len(hiddenOptionTypes) > 0 {
		hiddenData := strings.Join(hiddenOptionTypes, "\n")
		err = os.WriteFile("generated/hiddenOptionTypes.tf", []byte(hiddenData), 0644)
		if err != nil {
			log.Println(err)
		}
	}

	if len(passwordOptionTypes) > 0 {
		passwordData := strings.Join(passwordOptionTypes, "\n")
		err = os.WriteFile("generated/passwordOptionTypes.tf", []byte(passwordData), 0644)
		if err != nil {
			log.Println(err)
		}
	}

	if len(numberOptionTypes) > 0 {
		numberData := strings.Join(numberOptionTypes, "\n")
		err = os.WriteFile("generated/numberOptionTypes.tf", []byte(numberData), 0644)
		if err != nil {
			log.Println(err)
		}
	}

	if len(typeAheadOptionTypes) > 0 {
		typeAheadData := strings.Join(typeAheadOptionTypes, "\n")
		err = os.WriteFile("generated/typeAheadOptionTypes.tf", []byte(typeAheadData), 0644)
		if err != nil {
			log.Println(err)
		}
	}

	if len(checkboxOptionTypes) > 0 {
		checkboxData := strings.Join(checkboxOptionTypes, "\n")
		err = os.WriteFile("generated/checkboxOptionTypes.tf", []byte(checkboxData), 0644)
		if err != nil {
			log.Println(err)
		}
	}
}

func generateTextOptionTypes(resource morpheus.OptionType) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_text_option_type", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("description", cty.StringVal(resource.Description))
	providerBody.SetAttributeValue("field_name", cty.StringVal(resource.FieldName))
	providerBody.SetAttributeValue("export_meta", cty.BoolVal(resource.ExportMeta))
	providerBody.SetAttributeValue("dependent_field", cty.StringVal(resource.DependsOnCode))
	providerBody.SetAttributeValue("visibility_field", cty.StringVal(resource.VisibleOnCode))
	providerBody.SetAttributeValue("display_value_on_details", cty.BoolVal(resource.DisplayValueOnDetails))
	providerBody.SetAttributeValue("field_label", cty.StringVal(resource.FieldLabel))
	providerBody.SetAttributeValue("placeholder", cty.StringVal(resource.PlaceHolder))
	providerBody.SetAttributeValue("default_value", cty.StringVal(resource.DefaultValue))
	providerBody.SetAttributeValue("help_block", cty.StringVal(resource.HelpBlock))
	providerBody.SetAttributeValue("required", cty.BoolVal(resource.Required))
	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateHiddenOptionTypes(resource morpheus.OptionType) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_hidden_option_type", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("description", cty.StringVal(resource.Description))
	providerBody.SetAttributeValue("field_name", cty.StringVal(resource.FieldName))
	providerBody.SetAttributeValue("export_meta", cty.BoolVal(resource.ExportMeta))
	providerBody.SetAttributeValue("dependent_field", cty.StringVal(resource.DependsOnCode))
	providerBody.SetAttributeValue("visibility_field", cty.StringVal(resource.VisibleOnCode))
	providerBody.SetAttributeValue("display_value_on_details", cty.BoolVal(resource.DisplayValueOnDetails))
	providerBody.SetAttributeValue("default_value", cty.StringVal(resource.DefaultValue))
	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateNumberOptionTypes(resource morpheus.OptionType) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_number_option_type", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("description", cty.StringVal(resource.Description))
	providerBody.SetAttributeValue("field_name", cty.StringVal(resource.FieldName))
	providerBody.SetAttributeValue("export_meta", cty.BoolVal(resource.ExportMeta))
	providerBody.SetAttributeValue("dependent_field", cty.StringVal(resource.DependsOnCode))
	providerBody.SetAttributeValue("visibility_field", cty.StringVal(resource.VisibleOnCode))
	providerBody.SetAttributeValue("display_value_on_details", cty.BoolVal(resource.DisplayValueOnDetails))
	providerBody.SetAttributeValue("field_label", cty.StringVal(resource.FieldLabel))
	providerBody.SetAttributeValue("placeholder", cty.StringVal(resource.PlaceHolder))
	providerBody.SetAttributeValue("default_value", cty.StringVal(resource.DefaultValue))
	providerBody.SetAttributeValue("help_block", cty.StringVal(resource.HelpBlock))
	providerBody.SetAttributeValue("required", cty.BoolVal(resource.Required))
	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generatePasswordOptionTypes(resource morpheus.OptionType) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_password_option_type", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("description", cty.StringVal(resource.Description))
	providerBody.SetAttributeValue("field_name", cty.StringVal(resource.FieldName))
	providerBody.SetAttributeValue("export_meta", cty.BoolVal(resource.ExportMeta))
	providerBody.SetAttributeValue("dependent_field", cty.StringVal(resource.DependsOnCode))
	providerBody.SetAttributeValue("visibility_field", cty.StringVal(resource.VisibleOnCode))
	providerBody.SetAttributeValue("display_value_on_details", cty.BoolVal(resource.DisplayValueOnDetails))
	providerBody.SetAttributeValue("field_label", cty.StringVal(resource.FieldLabel))
	providerBody.SetAttributeValue("placeholder", cty.StringVal(resource.PlaceHolder))
	providerBody.SetAttributeValue("default_value", cty.StringVal(resource.DefaultValue))
	providerBody.SetAttributeValue("help_block", cty.StringVal(resource.HelpBlock))
	providerBody.SetAttributeValue("required", cty.BoolVal(resource.Required))
	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateSelecListOptionTypes(resource morpheus.OptionType) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_select_list_option_type", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("description", cty.StringVal(resource.Description))
	providerBody.SetAttributeValue("field_name", cty.StringVal(resource.FieldName))
	providerBody.SetAttributeValue("export_meta", cty.BoolVal(resource.ExportMeta))
	providerBody.SetAttributeValue("dependent_field", cty.StringVal(resource.DependsOnCode))
	providerBody.SetAttributeValue("visibility_field", cty.StringVal(resource.VisibleOnCode))
	providerBody.SetAttributeValue("display_value_on_details", cty.BoolVal(resource.DisplayValueOnDetails))
	providerBody.SetAttributeValue("field_label", cty.StringVal(resource.FieldLabel))
	providerBody.SetAttributeValue("default_value", cty.StringVal(resource.DefaultValue))
	providerBody.SetAttributeValue("help_block", cty.StringVal(resource.HelpBlock))
	providerBody.SetAttributeValue("required", cty.BoolVal(resource.Required))

	if resource.OptionListId != nil {
		optionID := resource.OptionListId
		oloutput := fmt.Sprintf("%v", optionID["id"])
		numout, err := cty.ParseNumberVal(oloutput)
		if err != nil {
			log.Println(err)
		}
		providerBody.SetAttributeValue("option_list_id", numout)
	}

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateTypeAheadOptionTypes(resource morpheus.OptionType) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_typeahead_option_type", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("description", cty.StringVal(resource.Description))
	providerBody.SetAttributeValue("field_name", cty.StringVal(resource.FieldName))
	providerBody.SetAttributeValue("export_meta", cty.BoolVal(resource.ExportMeta))
	providerBody.SetAttributeValue("dependent_field", cty.StringVal(resource.DependsOnCode))
	providerBody.SetAttributeValue("visibility_field", cty.StringVal(resource.VisibleOnCode))
	providerBody.SetAttributeValue("display_value_on_details", cty.BoolVal(resource.DisplayValueOnDetails))
	providerBody.SetAttributeValue("field_label", cty.StringVal(resource.FieldLabel))
	providerBody.SetAttributeValue("default_value", cty.StringVal(resource.DefaultValue))
	providerBody.SetAttributeValue("help_block", cty.StringVal(resource.HelpBlock))
	providerBody.SetAttributeValue("required", cty.BoolVal(resource.Required))
	optionID := resource.OptionListId
	oloutput := fmt.Sprintf("%v", optionID["id"])
	numout, err := cty.ParseNumberVal(oloutput)
	if err != nil {
		log.Println(err)
	}
	providerBody.SetAttributeValue("option_list_id", numout)

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateCheckboxOptionTypes(resource morpheus.OptionType) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_checkbox_option_type", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("description", cty.StringVal(resource.Description))
	providerBody.SetAttributeValue("field_name", cty.StringVal(resource.FieldName))
	providerBody.SetAttributeValue("export_meta", cty.BoolVal(resource.ExportMeta))
	providerBody.SetAttributeValue("dependent_field", cty.StringVal(resource.DependsOnCode))
	providerBody.SetAttributeValue("visibility_field", cty.StringVal(resource.VisibleOnCode))
	providerBody.SetAttributeValue("display_value_on_details", cty.BoolVal(resource.DisplayValueOnDetails))
	providerBody.SetAttributeValue("field_label", cty.StringVal(resource.FieldLabel))
	providerBody.SetAttributeValue("default_checked", cty.StringVal(resource.DefaultValue))
	hcloutput := string(hclFile.Bytes())
	return hcloutput
}
