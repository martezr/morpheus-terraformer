package resources

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/hashicorp/go-getter/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func InstallProvider() {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	providerDir := dirname + "/.terraform.d/plugins/morpheusdata.com/gomorpheus/morpheus/0.3.1/darwin_amd64"
	if _, err := os.Stat(providerDir); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(providerDir, 0755)
		check(err)
	}
	client := &getter.Client{}
	request := &getter.Request{
		Src: "https://github.com/gomorpheus/terraform-provider-morpheus/releases/download/v0.3.1/terraform-provider-morpheus_0.3.1_darwin_amd64.zip",
		Dst: providerDir,
	}
	_, err = client.Get(context.TODO(), request)
	if err != nil {
		log.Printf("Error getting path: %v", err)
	}

	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// create new file on system
	tfFile, err := os.Create("generated/provider.tf")
	if err != nil {
		log.Println(err)
		return
	}
	// initialize the body of the new file object
	rootBody := hclFile.Body()
	provider := rootBody.AppendNewBlock("provider",
		[]string{"morpheus"})
	providerBody := provider.Body()
	providerBody.SetAttributeValue("url", cty.StringVal("https://morpheus.test.local"))

	tfBlock := rootBody.AppendNewBlock("terraform", nil)
	tfBlockBody := tfBlock.Body()
	reqProvsBlock := tfBlockBody.AppendNewBlock("required_providers",
		nil)
	reqProvsBlockBody := reqProvsBlock.Body()
	reqProvsBlockBody.SetAttributeValue("morpheus",
		cty.ObjectVal(map[string]cty.Value{
			"source":  cty.StringVal("morpheusdata.com/gomorpheus/morpheus"),
			"version": cty.StringVal("0.3.1"),
		}))
	tfFile.Write(hclFile.Bytes())

}

/*
func runTerraform() {
	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("1.0.6")),
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		log.Fatalf("error installing Terraform: %s", err)
	}

	workingDir := "./generated"
	tf, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	}
	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Fatalf("error running Init: %s", err)
	}

	state, err := tf.Show(context.Background())
	if err != nil {
		log.Fatalf("error running Show: %s", err)
	}

	planResonse, err := tf.Plan(context.Background())
	if err != nil {
		log.Fatalf("error running Plan: %s", err)
	}

	//log.Println(tf.Validate(context.Background()))

	tf.FormatWrite(context.Background())

	log.Println(planResonse)

	log.Println(state.FormatVersion) // "0.1"
}
*/
func check(e error) {
	if e != nil {
		panic(e)
	}
}
