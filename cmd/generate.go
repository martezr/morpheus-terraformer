package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/martezr/morpheus-terraformer/resources"
	"github.com/martezr/morpheus-terraformer/utils"
	"github.com/spf13/cobra"
)

var genResources []string
var filterResources []string
var excludeResources []string

func init() {
	generateCmd.Flags().StringSliceVarP(&genResources, "resources", "r", []string{}, "groups,environments or * for all services")
	generateCmd.MarkFlagRequired("resources")
	generateCmd.Flags().StringSliceVarP(&filterResources, "filter", "f", []string{}, "filter the resources returned")
	generateCmd.Flags().StringSliceVarP(&excludeResources, "exclude", "e", []string{}, "exclude specific resources when used with '*' to return a subset of all the resources")
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Terraform code from existing Morpheus resources",
	Run: func(cmd *cobra.Command, args []string) {
		baseDir := "generated"
		if _, err := os.Stat(baseDir); errors.Is(err, os.ErrNotExist) {
			err = os.MkdirAll(baseDir, 0755)
			if err != nil {
				panic(err)
			}
		}

		morpheusURL := os.Getenv("MORPHEUS_API_URL")
		morpheusToken := os.Getenv("MORPHEUS_API_TOKEN")

		morpheusUsername := os.Getenv("MORPHEUS_API_USERNAME")
		morpheusPassword := os.Getenv("MORPHEUS_API_PASSWORD")
		var client *morpheus.Client

		if morpheusURL != "" {
			client = morpheus.NewClient(morpheusURL)
		} else {
			fmt.Println("The MORPHEUS_API_URL is not specified")
			os.Exit(1)
		}

		if morpheusToken != "" {
			client.SetAccessToken("", "", 0, "write")
		} else if morpheusUsername != "" && morpheusPassword != "" {
			client.SetUsernameAndPassword(morpheusUsername, morpheusPassword)
		} else {
			fmt.Println("Login credentials not specified")
			os.Exit(1)
		}

		log.Println(genResources)
		log.Println(filterResources)
		log.Println(excludeResources)

		if utils.Contains(genResources, "*") && len(filterResources) >= 0 {
			if !utils.Contains(excludeResources, "environment") {
				resources.GenerateEnvironments(client)
			}
			if !utils.Contains(excludeResources, "group") {
				resources.GenerateGroups(client)
			}
			if !utils.Contains(excludeResources, "task") {
				resources.GenerateTasks(client)
			}

			if !utils.Contains(excludeResources, "tenant") {
				resources.GenerateTenants(client)
			}

			if !utils.Contains(excludeResources, "optiontype") {
				resources.GenerateOptionTypes(client)
			}
			if !utils.Contains(excludeResources, "optionlist") {
				resources.GenerateOptionLists(client)
			}

			if !utils.Contains(excludeResources, "optiontype") {
				resources.GenerateWorkflows(client)
			}
		} else {
			for _, resource := range genResources {
				switch resource {
				case "group":
					resources.GenerateGroups(client)
				case "environment":
					resources.GenerateEnvironments(client)
				case "optiontype":
					resources.GenerateOptionTypes(client)
				case "optionlist":
					resources.GenerateOptionLists(client)
				case "task":
					resources.GenerateTasks(client)
				case "tenant":
					resources.GenerateTenants(client)
				case "workflow":
					resources.GenerateWorkflows(client)
				default:
					log.Printf("unable to generate resources for %s as %s is an invalid resource type", resource, resource)
				}
			}
		}
	},
}
