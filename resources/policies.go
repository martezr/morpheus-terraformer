package resources

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/martezr/morpheus-terraformer/utils"
	"github.com/zclconf/go-cty/cty"
)

// GenerateOptionLists generates terraform code for Morpheus option lists
func GeneratePolicies(client *morpheus.Client) {
	log.Println("generating policies...")
	response, err := client.ListPolicies(&morpheus.Request{
		QueryParams: map[string]string{"max": "500"},
	})
	if err != nil {
		log.Println(err)
	}

	result := response.Result.(*morpheus.ListPoliciesResult)
	policies := result.Policies
	var backupCreationPolicies []string
	var budgetPolicies []string
	var clusterNamePolicies []string
	var hostnamePolicies []string
	var instanceNamePolicies []string
	var maxContainersPolicies []string
	var maxCoresPolicies []string
	var maxHostsPolicies []string
	var maxMemoryPolicies []string
	var maxStoragePolicies []string
	var maxVmsPolicies []string
	var networkQuotaPolicies []string
	var routerQuotaPolicies []string
	var userCreationPolicies []string
	var workflowPolicies []string

	for _, v := range *policies {
		switch v.PolicyType.Code {
		case "createBackup":
			backupCreationPolicies = append(backupCreationPolicies, generateBackupCreationPolicy(v))
		case "createUser":
			userCreationPolicies = append(userCreationPolicies, generateUserCreationPolicy(v))
		case "createUserGroup":
		case "delayedRemoval":
		case "lifecycle":
		case "hostNaming":
			hostnamePolicies = append(hostnamePolicies, generateHostNamePolicy(v))
		case "maxContainers":
			maxContainersPolicies = append(maxContainersPolicies, generateMaxContainersPolicy(v))
		case "maxCores":
			maxCoresPolicies = append(maxCoresPolicies, generateMaxCoresPolicy(v))
		case "maxHosts":
			maxHostsPolicies = append(maxHostsPolicies, generateMaxHostsPolicy(v))
		case "maxMemory":
			maxMemoryPolicies = append(maxMemoryPolicies, generateMaxMemoryPolicy(v))
		case "maxNetworks":
			networkQuotaPolicies = append(networkQuotaPolicies, generateNetworkQuotaPolicy(v))
		// Budget Policy
		case "maxPrice":
			budgetPolicies = append(budgetPolicies, generateBudgetPolicy(v))
		case "maxRouters":
			routerQuotaPolicies = append(routerQuotaPolicies, generateRouterQuotaPolicy(v))
		case "maxStorage":
			maxStoragePolicies = append(maxStoragePolicies, generateMaxStoragePolicy(v))
		case "maxVms":
			maxVmsPolicies = append(maxVmsPolicies, generateMaxVmsPolicy(v))
		case "motd":
		// Instance Naming Policy
		case "naming":
			instanceNamePolicies = append(instanceNamePolicies, generateInstanceNamePolicy(v))
		case "powerSchedule":
		case "provisionApproval":
		case "reconfigureApproval":
		// Cluster Resource Naming Policy
		case "serverNaming":
			clusterNamePolicies = append(clusterNamePolicies, generateClusterNamePolicy(v))
		case "shutdown":
		case "tags":
		case "workflow":
			workflowPolicies = append(workflowPolicies, generateWorkflowPolicy(v))
		}
	}

	// Write backup creation policies
	backupCreationPolicyData := strings.Join(backupCreationPolicies, "\n")
	err = os.WriteFile("generated/backupCreationPolicies.tf", []byte(backupCreationPolicyData), 0644)
	if err != nil {
		log.Println(err)
	}

	// Write budget policies
	budgetPolicyData := strings.Join(budgetPolicies, "\n")
	err = os.WriteFile("generated/budgetPolicies.tf", []byte(budgetPolicyData), 0644)
	if err != nil {
		log.Println(err)
	}

	// Write cluster name policies
	clusterNamePolicyData := strings.Join(clusterNamePolicies, "\n")
	err = os.WriteFile("generated/clusterNamePolicies.tf", []byte(clusterNamePolicyData), 0644)
	if err != nil {
		log.Println(err)
	}

	// Write host name policies
	hostNamePolicyData := strings.Join(hostnamePolicies, "\n")
	err = os.WriteFile("generated/hostNamePolicies.tf", []byte(hostNamePolicyData), 0644)
	if err != nil {
		log.Println(err)
	}

	// Write instance name policies
	instanceNamePolicyData := strings.Join(instanceNamePolicies, "\n")
	err = os.WriteFile("generated/instanceNamePolicies.tf", []byte(instanceNamePolicyData), 0644)
	if err != nil {
		log.Println(err)
	}

	// Write max containers policies
	maxContainersPolicyData := strings.Join(maxContainersPolicies, "\n")
	err = os.WriteFile("generated/maxContainersPolicies.tf", []byte(maxContainersPolicyData), 0644)
	if err != nil {
		log.Println(err)
	}

	// Write max cores policies
	maxCoresPolicyData := strings.Join(maxCoresPolicies, "\n")
	err = os.WriteFile("generated/maxCoresPolicies.tf", []byte(maxCoresPolicyData), 0644)
	if err != nil {
		log.Println(err)
	}

	// Write max hosts policies
	maxHostsPolicyData := strings.Join(maxHostsPolicies, "\n")
	err = os.WriteFile("generated/maxHostsPolicies.tf", []byte(maxHostsPolicyData), 0644)
	if err != nil {
		log.Println(err)
	}

	// Write max memory policies
	maxMemoryPolicyData := strings.Join(maxMemoryPolicies, "\n")
	err = os.WriteFile("generated/maxMemoryPolicies.tf", []byte(maxMemoryPolicyData), 0644)
	if err != nil {
		log.Println(err)
	}

	// Write max storage policies
	maxStoragePolicyData := strings.Join(maxStoragePolicies, "\n")
	err = os.WriteFile("generated/maxStoragePolicies.tf", []byte(maxStoragePolicyData), 0644)
	if err != nil {
		log.Println(err)
	}

	// Write max vms policies
	maxVmsPolicyData := strings.Join(maxVmsPolicies, "\n")
	err = os.WriteFile("generated/maxVmsPolicies.tf", []byte(maxVmsPolicyData), 0644)
	if err != nil {
		log.Println(err)
	}

	// Write network quota policies
	networkQuotaPolicyData := strings.Join(networkQuotaPolicies, "\n")
	err = os.WriteFile("generated/networkQuotaPolicies.tf", []byte(networkQuotaPolicyData), 0644)
	if err != nil {
		log.Println(err)
	}

	// Write router quota policies
	routerQuotaPolicyData := strings.Join(routerQuotaPolicies, "\n")
	err = os.WriteFile("generated/routerQuotaPolicies.tf", []byte(routerQuotaPolicyData), 0644)
	if err != nil {
		log.Println(err)
	}

	// Write user creation policies
	userCreationPolicyData := strings.Join(userCreationPolicies, "\n")
	err = os.WriteFile("generated/userCreationPolicies.tf", []byte(userCreationPolicyData), 0644)
	if err != nil {
		log.Println(err)
	}

	// Write workflow policies
	v := strings.Join(workflowPolicies, "\n")
	err = os.WriteFile("generated/workflowPolicies.tf", []byte(v), 0644)
	if err != nil {
		log.Println(err)
	}
}

func generateBackupCreationPolicy(resource morpheus.Policy) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_backup_creation_policy", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("description", cty.StringVal(resource.Description))
	providerBody.SetAttributeValue("enabled", cty.BoolVal(resource.Enabled))
	providerBody.SetAttributeValue("enforcement_type", cty.StringVal(resource.Config.CreateBackupType))
	if resource.Config.CreateBackup == "on" {
		providerBody.SetAttributeValue("create_backup", cty.BoolVal(true))
	} else {
		providerBody.SetAttributeValue("create_backup", cty.BoolVal(false))
	}

	switch resource.RefType {
	case "ComputeSite":
		providerBody.SetAttributeValue("scope", cty.StringVal("group"))
		providerBody.SetAttributeValue("group_id", cty.NumberIntVal(resource.Site.ID))
	case "ComputeZone":
		providerBody.SetAttributeValue("scope", cty.StringVal("cloud"))
		providerBody.SetAttributeValue("cloud_id", cty.NumberIntVal(resource.Zone.ID))
	case "User":
		providerBody.SetAttributeValue("scope", cty.StringVal("user"))
		providerBody.SetAttributeValue("user_id", cty.NumberIntVal(resource.User.ID))
	case "Role":
		providerBody.SetAttributeValue("scope", cty.StringVal("role"))
		providerBody.SetAttributeValue("role_id", cty.NumberIntVal(resource.Role.ID))
		providerBody.SetAttributeValue("apply_to_each_user", cty.BoolVal(resource.EachUser))
	default:
		providerBody.SetAttributeValue("scope", cty.StringVal("global"))
	}

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateBudgetPolicy(resource morpheus.Policy) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_budget_policy", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("description", cty.StringVal(resource.Description))
	providerBody.SetAttributeValue("enabled", cty.BoolVal(resource.Enabled))
	providerBody.SetAttributeValue("max_price", cty.StringVal(resource.Config.MaxPrice))
	providerBody.SetAttributeValue("currency", cty.StringVal(resource.Config.MaxPriceCurrency))
	providerBody.SetAttributeValue("unit_of_time", cty.StringVal(resource.Config.MaxPriceUnit))

	switch resource.RefType {
	case "ComputeSite":
		providerBody.SetAttributeValue("scope", cty.StringVal("group"))
		providerBody.SetAttributeValue("group_id", cty.NumberIntVal(resource.Site.ID))
	case "ComputeZone":
		providerBody.SetAttributeValue("scope", cty.StringVal("cloud"))
		providerBody.SetAttributeValue("cloud_id", cty.NumberIntVal(resource.Zone.ID))
	case "User":
		providerBody.SetAttributeValue("scope", cty.StringVal("user"))
		providerBody.SetAttributeValue("user_id", cty.NumberIntVal(resource.User.ID))
	case "Role":
		providerBody.SetAttributeValue("scope", cty.StringVal("role"))
		providerBody.SetAttributeValue("role_id", cty.NumberIntVal(resource.Role.ID))
		providerBody.SetAttributeValue("apply_to_each_user", cty.BoolVal(resource.EachUser))
	default:
		providerBody.SetAttributeValue("scope", cty.StringVal("global"))
	}

	var tenant_ids []cty.Value
	for _, account := range resource.Accounts {
		tenant_ids = append(tenant_ids, cty.NumberIntVal(account.ID))
	}
	if len(tenant_ids) == 0 {
		tenantId := cty.String
		providerBody.SetAttributeValue("tenant_ids", cty.ListValEmpty(tenantId))
	} else {
		providerBody.SetAttributeValue("tenant_ids", cty.ListVal(tenant_ids))
	}
	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateClusterNamePolicy(resource morpheus.Policy) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_cluster_resource_name_policy", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("description", cty.StringVal(resource.Description))
	providerBody.SetAttributeValue("enabled", cty.BoolVal(resource.Enabled))
	providerBody.SetAttributeValue("enforcement_type", cty.StringVal(resource.Config.ServerNamingType))
	providerBody.SetAttributeValue("naming_pattern", cty.StringVal(resource.Config.ServerNamingPattern))
	if resource.Config.ServerNamingConflict == "on" {
		providerBody.SetAttributeValue("auto_resolve_conflicts", cty.BoolVal(true))
	} else {
		providerBody.SetAttributeValue("auto_resolve_conflicts", cty.BoolVal(false))
	}

	switch resource.RefType {
	case "ComputeSite":
		providerBody.SetAttributeValue("scope", cty.StringVal("group"))
		providerBody.SetAttributeValue("group_id", cty.NumberIntVal(resource.Site.ID))
	case "ComputeZone":
		providerBody.SetAttributeValue("scope", cty.StringVal("cloud"))
		providerBody.SetAttributeValue("cloud_id", cty.NumberIntVal(resource.Zone.ID))
	case "User":
		providerBody.SetAttributeValue("scope", cty.StringVal("user"))
		providerBody.SetAttributeValue("user_id", cty.NumberIntVal(resource.User.ID))
	case "Role":
		providerBody.SetAttributeValue("scope", cty.StringVal("role"))
		providerBody.SetAttributeValue("role_id", cty.NumberIntVal(resource.Role.ID))
		providerBody.SetAttributeValue("apply_to_each_user", cty.BoolVal(resource.EachUser))
	default:
		providerBody.SetAttributeValue("scope", cty.StringVal("global"))
	}

	var tenant_ids []cty.Value
	for _, account := range resource.Accounts {
		tenant_ids = append(tenant_ids, cty.NumberIntVal(account.ID))
	}
	if len(tenant_ids) == 0 {
		tenantId := cty.String
		providerBody.SetAttributeValue("tenant_ids", cty.ListValEmpty(tenantId))
	} else {
		providerBody.SetAttributeValue("tenant_ids", cty.ListVal(tenant_ids))
	}

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateHostNamePolicy(resource morpheus.Policy) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_hostname_policy", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("description", cty.StringVal(resource.Description))
	providerBody.SetAttributeValue("enabled", cty.BoolVal(resource.Enabled))
	providerBody.SetAttributeValue("enforcement_type", cty.StringVal(resource.Config.HostNamingType))
	providerBody.SetAttributeValue("naming_pattern", cty.StringVal(resource.Config.HostNamingPattern))

	switch resource.RefType {
	case "ComputeSite":
		providerBody.SetAttributeValue("scope", cty.StringVal("group"))
		providerBody.SetAttributeValue("group_id", cty.NumberIntVal(resource.Site.ID))
	case "ComputeZone":
		providerBody.SetAttributeValue("scope", cty.StringVal("cloud"))
		providerBody.SetAttributeValue("cloud_id", cty.NumberIntVal(resource.Zone.ID))
	case "User":
		providerBody.SetAttributeValue("scope", cty.StringVal("user"))
		providerBody.SetAttributeValue("user_id", cty.NumberIntVal(resource.User.ID))
	case "Role":
		providerBody.SetAttributeValue("scope", cty.StringVal("role"))
		providerBody.SetAttributeValue("role_id", cty.NumberIntVal(resource.Role.ID))
		providerBody.SetAttributeValue("apply_to_each_user", cty.BoolVal(resource.EachUser))
	default:
		providerBody.SetAttributeValue("scope", cty.StringVal("global"))
	}

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateInstanceNamePolicy(resource morpheus.Policy) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_instance_name_policy", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("description", cty.StringVal(resource.Description))
	providerBody.SetAttributeValue("enabled", cty.BoolVal(resource.Enabled))
	providerBody.SetAttributeValue("enforcement_type", cty.StringVal(resource.Config.NamingType))
	providerBody.SetAttributeValue("naming_pattern", cty.StringVal(resource.Config.NamingPattern))
	if resource.Config.NamingConflict == "on" {
		providerBody.SetAttributeValue("auto_resolve_conflicts", cty.BoolVal(true))
	} else {
		providerBody.SetAttributeValue("auto_resolve_conflicts", cty.BoolVal(false))
	}

	switch resource.RefType {
	case "ComputeSite":
		providerBody.SetAttributeValue("scope", cty.StringVal("group"))
		providerBody.SetAttributeValue("group_id", cty.NumberIntVal(resource.Site.ID))
	case "ComputeZone":
		providerBody.SetAttributeValue("scope", cty.StringVal("cloud"))
		providerBody.SetAttributeValue("cloud_id", cty.NumberIntVal(resource.Zone.ID))
	case "User":
		providerBody.SetAttributeValue("scope", cty.StringVal("user"))
		providerBody.SetAttributeValue("user_id", cty.NumberIntVal(resource.User.ID))
	case "Role":
		providerBody.SetAttributeValue("scope", cty.StringVal("role"))
		providerBody.SetAttributeValue("role_id", cty.NumberIntVal(resource.Role.ID))
		providerBody.SetAttributeValue("apply_to_each_user", cty.BoolVal(resource.EachUser))
	default:
		providerBody.SetAttributeValue("scope", cty.StringVal("global"))
	}

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateMaxContainersPolicy(resource morpheus.Policy) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_max_containers_policy", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("description", cty.StringVal(resource.Description))
	providerBody.SetAttributeValue("enabled", cty.BoolVal(resource.Enabled))

	n, err := strconv.Atoi(resource.Config.MaxContainers)
	if err != nil {
		log.Println(err)
	}
	providerBody.SetAttributeValue("max_containers", cty.NumberIntVal(int64(n)))

	switch resource.RefType {
	case "ComputeSite":
		providerBody.SetAttributeValue("scope", cty.StringVal("group"))
		providerBody.SetAttributeValue("group_id", cty.NumberIntVal(resource.Site.ID))
	case "ComputeZone":
		providerBody.SetAttributeValue("scope", cty.StringVal("cloud"))
		providerBody.SetAttributeValue("cloud_id", cty.NumberIntVal(resource.Zone.ID))
	case "User":
		providerBody.SetAttributeValue("scope", cty.StringVal("user"))
		providerBody.SetAttributeValue("user_id", cty.NumberIntVal(resource.User.ID))
	case "Role":
		providerBody.SetAttributeValue("scope", cty.StringVal("role"))
		providerBody.SetAttributeValue("role_id", cty.NumberIntVal(resource.Role.ID))
		providerBody.SetAttributeValue("apply_to_each_user", cty.BoolVal(resource.EachUser))
	default:
		providerBody.SetAttributeValue("scope", cty.StringVal("global"))
	}

	var tenant_ids []cty.Value
	for _, account := range resource.Accounts {
		tenant_ids = append(tenant_ids, cty.NumberIntVal(account.ID))
	}
	if len(tenant_ids) == 0 {
		tenantId := cty.String
		providerBody.SetAttributeValue("tenant_ids", cty.ListValEmpty(tenantId))
	} else {
		providerBody.SetAttributeValue("tenant_ids", cty.ListVal(tenant_ids))
	}

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateMaxCoresPolicy(resource morpheus.Policy) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_max_cores_policy", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("description", cty.StringVal(resource.Description))
	providerBody.SetAttributeValue("enabled", cty.BoolVal(resource.Enabled))

	n, err := strconv.Atoi(resource.Config.MaxCores)
	if err != nil {
		log.Println(err)
	}
	providerBody.SetAttributeValue("max_cores", cty.NumberIntVal(int64(n)))

	switch resource.RefType {
	case "ComputeSite":
		providerBody.SetAttributeValue("scope", cty.StringVal("group"))
		providerBody.SetAttributeValue("group_id", cty.NumberIntVal(resource.Site.ID))
	case "ComputeZone":
		providerBody.SetAttributeValue("scope", cty.StringVal("cloud"))
		providerBody.SetAttributeValue("cloud_id", cty.NumberIntVal(resource.Zone.ID))
	case "User":
		providerBody.SetAttributeValue("scope", cty.StringVal("user"))
		providerBody.SetAttributeValue("user_id", cty.NumberIntVal(resource.User.ID))
	case "Role":
		providerBody.SetAttributeValue("scope", cty.StringVal("role"))
		providerBody.SetAttributeValue("role_id", cty.NumberIntVal(resource.Role.ID))
		providerBody.SetAttributeValue("apply_to_each_user", cty.BoolVal(resource.EachUser))
	default:
		providerBody.SetAttributeValue("scope", cty.StringVal("global"))
	}

	var tenant_ids []cty.Value
	for _, account := range resource.Accounts {
		tenant_ids = append(tenant_ids, cty.NumberIntVal(account.ID))
	}
	if len(tenant_ids) == 0 {
		tenantId := cty.String
		providerBody.SetAttributeValue("tenant_ids", cty.ListValEmpty(tenantId))
	} else {
		providerBody.SetAttributeValue("tenant_ids", cty.ListVal(tenant_ids))
	}

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateMaxHostsPolicy(resource morpheus.Policy) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_max_hosts_policy", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("description", cty.StringVal(resource.Description))
	providerBody.SetAttributeValue("enabled", cty.BoolVal(resource.Enabled))

	n, err := strconv.Atoi(resource.Config.MaxHosts)
	if err != nil {
		log.Println(err)
	}
	providerBody.SetAttributeValue("max_hosts", cty.NumberIntVal(int64(n)))

	switch resource.RefType {
	case "ComputeSite":
		providerBody.SetAttributeValue("scope", cty.StringVal("group"))
		providerBody.SetAttributeValue("group_id", cty.NumberIntVal(resource.Site.ID))
	case "ComputeZone":
		providerBody.SetAttributeValue("scope", cty.StringVal("cloud"))
		providerBody.SetAttributeValue("cloud_id", cty.NumberIntVal(resource.Zone.ID))
	case "User":
		providerBody.SetAttributeValue("scope", cty.StringVal("user"))
		providerBody.SetAttributeValue("user_id", cty.NumberIntVal(resource.User.ID))
	case "Role":
		providerBody.SetAttributeValue("scope", cty.StringVal("role"))
		providerBody.SetAttributeValue("role_id", cty.NumberIntVal(resource.Role.ID))
		providerBody.SetAttributeValue("apply_to_each_user", cty.BoolVal(resource.EachUser))
	default:
		providerBody.SetAttributeValue("scope", cty.StringVal("global"))
	}

	var tenant_ids []cty.Value
	for _, account := range resource.Accounts {
		tenant_ids = append(tenant_ids, cty.NumberIntVal(account.ID))
	}
	if len(tenant_ids) == 0 {
		tenantId := cty.String
		providerBody.SetAttributeValue("tenant_ids", cty.ListValEmpty(tenantId))
	} else {
		providerBody.SetAttributeValue("tenant_ids", cty.ListVal(tenant_ids))
	}

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateMaxMemoryPolicy(resource morpheus.Policy) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_max_memory_policy", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("description", cty.StringVal(resource.Description))
	providerBody.SetAttributeValue("enabled", cty.BoolVal(resource.Enabled))

	n, err := strconv.Atoi(resource.Config.MaxMemory)
	if err != nil {
		log.Println(err)
	}
	providerBody.SetAttributeValue("max_memory", cty.NumberIntVal(int64(n)))

	switch resource.RefType {
	case "ComputeSite":
		providerBody.SetAttributeValue("scope", cty.StringVal("group"))
		providerBody.SetAttributeValue("group_id", cty.NumberIntVal(resource.Site.ID))
	case "ComputeZone":
		providerBody.SetAttributeValue("scope", cty.StringVal("cloud"))
		providerBody.SetAttributeValue("cloud_id", cty.NumberIntVal(resource.Zone.ID))
	case "User":
		providerBody.SetAttributeValue("scope", cty.StringVal("user"))
		providerBody.SetAttributeValue("user_id", cty.NumberIntVal(resource.User.ID))
	case "Role":
		providerBody.SetAttributeValue("scope", cty.StringVal("role"))
		providerBody.SetAttributeValue("role_id", cty.NumberIntVal(resource.Role.ID))
		providerBody.SetAttributeValue("apply_to_each_user", cty.BoolVal(resource.EachUser))
	default:
		providerBody.SetAttributeValue("scope", cty.StringVal("global"))
	}

	var tenant_ids []cty.Value
	for _, account := range resource.Accounts {
		tenant_ids = append(tenant_ids, cty.NumberIntVal(account.ID))
	}
	if len(tenant_ids) == 0 {
		tenantId := cty.String
		providerBody.SetAttributeValue("tenant_ids", cty.ListValEmpty(tenantId))
	} else {
		providerBody.SetAttributeValue("tenant_ids", cty.ListVal(tenant_ids))
	}

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateMaxStoragePolicy(resource morpheus.Policy) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_max_storage_policy", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("description", cty.StringVal(resource.Description))
	providerBody.SetAttributeValue("enabled", cty.BoolVal(resource.Enabled))

	n, err := strconv.Atoi(resource.Config.MaxStorage)
	if err != nil {
		log.Println(err)
	}
	providerBody.SetAttributeValue("max_storage", cty.NumberIntVal(int64(n)))

	switch resource.RefType {
	case "ComputeSite":
		providerBody.SetAttributeValue("scope", cty.StringVal("group"))
		providerBody.SetAttributeValue("group_id", cty.NumberIntVal(resource.Site.ID))
	case "ComputeZone":
		providerBody.SetAttributeValue("scope", cty.StringVal("cloud"))
		providerBody.SetAttributeValue("cloud_id", cty.NumberIntVal(resource.Zone.ID))
	case "User":
		providerBody.SetAttributeValue("scope", cty.StringVal("user"))
		providerBody.SetAttributeValue("user_id", cty.NumberIntVal(resource.User.ID))
	case "Role":
		providerBody.SetAttributeValue("scope", cty.StringVal("role"))
		providerBody.SetAttributeValue("role_id", cty.NumberIntVal(resource.Role.ID))
		providerBody.SetAttributeValue("apply_to_each_user", cty.BoolVal(resource.EachUser))
	default:
		providerBody.SetAttributeValue("scope", cty.StringVal("global"))
	}

	var tenant_ids []cty.Value
	for _, account := range resource.Accounts {
		tenant_ids = append(tenant_ids, cty.NumberIntVal(account.ID))
	}
	if len(tenant_ids) == 0 {
		tenantId := cty.String
		providerBody.SetAttributeValue("tenant_ids", cty.ListValEmpty(tenantId))
	} else {
		providerBody.SetAttributeValue("tenant_ids", cty.ListVal(tenant_ids))
	}

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateMaxVmsPolicy(resource morpheus.Policy) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_max_vms_policy", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("description", cty.StringVal(resource.Description))
	providerBody.SetAttributeValue("enabled", cty.BoolVal(resource.Enabled))

	n, err := strconv.Atoi(resource.Config.MaxVms)
	if err != nil {
		log.Println(err)
	}
	providerBody.SetAttributeValue("max_vms", cty.NumberIntVal(int64(n)))

	switch resource.RefType {
	case "ComputeSite":
		providerBody.SetAttributeValue("scope", cty.StringVal("group"))
		providerBody.SetAttributeValue("group_id", cty.NumberIntVal(resource.Site.ID))
	case "ComputeZone":
		providerBody.SetAttributeValue("scope", cty.StringVal("cloud"))
		providerBody.SetAttributeValue("cloud_id", cty.NumberIntVal(resource.Zone.ID))
	case "User":
		providerBody.SetAttributeValue("scope", cty.StringVal("user"))
		providerBody.SetAttributeValue("user_id", cty.NumberIntVal(resource.User.ID))
	case "Role":
		providerBody.SetAttributeValue("scope", cty.StringVal("role"))
		providerBody.SetAttributeValue("role_id", cty.NumberIntVal(resource.Role.ID))
		providerBody.SetAttributeValue("apply_to_each_user", cty.BoolVal(resource.EachUser))
	default:
		providerBody.SetAttributeValue("scope", cty.StringVal("global"))
	}

	var tenant_ids []cty.Value
	for _, account := range resource.Accounts {
		tenant_ids = append(tenant_ids, cty.NumberIntVal(account.ID))
	}
	if len(tenant_ids) == 0 {
		tenantId := cty.String
		providerBody.SetAttributeValue("tenant_ids", cty.ListValEmpty(tenantId))
	} else {
		providerBody.SetAttributeValue("tenant_ids", cty.ListVal(tenant_ids))
	}

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateNetworkQuotaPolicy(resource morpheus.Policy) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_network_quota_policy", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("description", cty.StringVal(resource.Description))
	providerBody.SetAttributeValue("enabled", cty.BoolVal(resource.Enabled))

	n, err := strconv.Atoi(resource.Config.MaxNetworks)
	if err != nil {
		log.Println(err)
	}
	providerBody.SetAttributeValue("max_networks", cty.NumberIntVal(int64(n)))

	switch resource.RefType {
	case "ComputeSite":
		providerBody.SetAttributeValue("scope", cty.StringVal("group"))
		providerBody.SetAttributeValue("group_id", cty.NumberIntVal(resource.Site.ID))
	case "ComputeZone":
		providerBody.SetAttributeValue("scope", cty.StringVal("cloud"))
		providerBody.SetAttributeValue("cloud_id", cty.NumberIntVal(resource.Zone.ID))
	case "User":
		providerBody.SetAttributeValue("scope", cty.StringVal("user"))
		providerBody.SetAttributeValue("user_id", cty.NumberIntVal(resource.User.ID))
	case "Role":
		providerBody.SetAttributeValue("scope", cty.StringVal("role"))
		providerBody.SetAttributeValue("role_id", cty.NumberIntVal(resource.Role.ID))
		providerBody.SetAttributeValue("apply_to_each_user", cty.BoolVal(resource.EachUser))
	default:
		providerBody.SetAttributeValue("scope", cty.StringVal("global"))
	}

	var tenant_ids []cty.Value
	for _, account := range resource.Accounts {
		tenant_ids = append(tenant_ids, cty.NumberIntVal(account.ID))
	}
	if len(tenant_ids) == 0 {
		tenantId := cty.String
		providerBody.SetAttributeValue("tenant_ids", cty.ListValEmpty(tenantId))
	} else {
		providerBody.SetAttributeValue("tenant_ids", cty.ListVal(tenant_ids))
	}

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateRouterQuotaPolicy(resource morpheus.Policy) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_router_quota_policy", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("description", cty.StringVal(resource.Description))
	providerBody.SetAttributeValue("enabled", cty.BoolVal(resource.Enabled))

	n, err := strconv.Atoi(resource.Config.MaxRouters)
	if err != nil {
		log.Println(err)
	}
	providerBody.SetAttributeValue("max_routers", cty.NumberIntVal(int64(n)))

	switch resource.RefType {
	case "ComputeSite":
		providerBody.SetAttributeValue("scope", cty.StringVal("group"))
		providerBody.SetAttributeValue("group_id", cty.NumberIntVal(resource.Site.ID))
	case "ComputeZone":
		providerBody.SetAttributeValue("scope", cty.StringVal("cloud"))
		providerBody.SetAttributeValue("cloud_id", cty.NumberIntVal(resource.Zone.ID))
	case "User":
		providerBody.SetAttributeValue("scope", cty.StringVal("user"))
		providerBody.SetAttributeValue("user_id", cty.NumberIntVal(resource.User.ID))
	case "Role":
		providerBody.SetAttributeValue("scope", cty.StringVal("role"))
		providerBody.SetAttributeValue("role_id", cty.NumberIntVal(resource.Role.ID))
		providerBody.SetAttributeValue("apply_to_each_user", cty.BoolVal(resource.EachUser))
	default:
		providerBody.SetAttributeValue("scope", cty.StringVal("global"))
	}

	var tenant_ids []cty.Value
	for _, account := range resource.Accounts {
		tenant_ids = append(tenant_ids, cty.NumberIntVal(account.ID))
	}
	if len(tenant_ids) == 0 {
		tenantId := cty.String
		providerBody.SetAttributeValue("tenant_ids", cty.ListValEmpty(tenantId))
	} else {
		providerBody.SetAttributeValue("tenant_ids", cty.ListVal(tenant_ids))
	}

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateUserCreationPolicy(resource morpheus.Policy) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_user_creation_policy", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("description", cty.StringVal(resource.Description))
	providerBody.SetAttributeValue("enabled", cty.BoolVal(resource.Enabled))
	providerBody.SetAttributeValue("enforcement_type", cty.StringVal(resource.Config.CreateUserType))
	if resource.Config.CreateUser == "on" {
		providerBody.SetAttributeValue("create_user", cty.BoolVal(true))
	} else {
		providerBody.SetAttributeValue("create_user", cty.BoolVal(false))
	}
	switch resource.RefType {
	case "ComputeSite":
		providerBody.SetAttributeValue("scope", cty.StringVal("group"))
		providerBody.SetAttributeValue("group_id", cty.NumberIntVal(resource.Site.ID))
	case "ComputeZone":
		providerBody.SetAttributeValue("scope", cty.StringVal("cloud"))
		providerBody.SetAttributeValue("cloud_id", cty.NumberIntVal(resource.Zone.ID))
	case "User":
		providerBody.SetAttributeValue("scope", cty.StringVal("user"))
		providerBody.SetAttributeValue("user_id", cty.NumberIntVal(resource.User.ID))
	case "Role":
		providerBody.SetAttributeValue("scope", cty.StringVal("role"))
		providerBody.SetAttributeValue("role_id", cty.NumberIntVal(resource.Role.ID))
		providerBody.SetAttributeValue("apply_to_each_user", cty.BoolVal(resource.EachUser))
	default:
		providerBody.SetAttributeValue("scope", cty.StringVal("global"))
	}

	var tenant_ids []cty.Value
	for _, account := range resource.Accounts {
		tenant_ids = append(tenant_ids, cty.NumberIntVal(account.ID))
	}
	if len(tenant_ids) == 0 {
		tenantId := cty.String
		providerBody.SetAttributeValue("tenant_ids", cty.ListValEmpty(tenantId))
	} else {
		providerBody.SetAttributeValue("tenant_ids", cty.ListVal(tenant_ids))
	}

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}

func generateWorkflowPolicy(resource morpheus.Policy) (output string) {
	// create new empty hcl file object
	hclFile := hclwrite.NewEmptyFile()

	// initialize the body of the new file object
	rootBody := hclFile.Body()
	title := utils.GenerateResourceName(resource.Name)
	provider := rootBody.AppendNewBlock("resource",
		[]string{"morpheus_workflow_policy", title})
	providerBody := provider.Body()

	providerBody.SetAttributeValue("name", cty.StringVal(resource.Name))
	providerBody.SetAttributeValue("description", cty.StringVal(resource.Description))
	providerBody.SetAttributeValue("enabled", cty.BoolVal(resource.Enabled))
	n, err := strconv.Atoi(resource.Config.WorkflowID.(string))
	if err != nil {
		log.Println(err)
	}
	providerBody.SetAttributeValue("workflow_id", cty.NumberIntVal(int64(n)))

	switch resource.RefType {
	case "ComputeSite":
		providerBody.SetAttributeValue("scope", cty.StringVal("group"))
		providerBody.SetAttributeValue("group_id", cty.NumberIntVal(resource.Site.ID))
	case "ComputeZone":
		providerBody.SetAttributeValue("scope", cty.StringVal("cloud"))
		providerBody.SetAttributeValue("cloud_id", cty.NumberIntVal(resource.Zone.ID))
	case "User":
		providerBody.SetAttributeValue("scope", cty.StringVal("user"))
		providerBody.SetAttributeValue("user_id", cty.NumberIntVal(resource.User.ID))
	case "Role":
		providerBody.SetAttributeValue("scope", cty.StringVal("role"))
		providerBody.SetAttributeValue("role_id", cty.NumberIntVal(resource.Role.ID))
		providerBody.SetAttributeValue("apply_to_each_user", cty.BoolVal(resource.EachUser))
	default:
		providerBody.SetAttributeValue("scope", cty.StringVal("global"))
	}

	hcloutput := string(hclFile.Bytes())
	return hcloutput
}
