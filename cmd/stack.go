// Copyright © 2019 Kenneth Pianeta
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/kpianeta/mytools/helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)
type cfnResource struct {
	ResourceID   string
	Type         string
	Stack        string
	Status       string
	ResourceName string
	LogicalName  string
}


var stackName *string

// stackCmd represents the stack command
var stackCmd = &cobra.Command{
	Use:   "stack",
	Short: "Parent command to query an AWS stack.",
	Long: `Parent command to query an AWS stack for resource data.
For example:


Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: listStackResources,
	//Run: func(cmd *cobra.Command, args []string) {
	//	stackName = cmd.PersistentFlags().StringP("name","n","", "AWS Stack name")
	//	fmt.Printf("parentStackResource called for parentStackResource: %s\n", cmd.PersistentFlags().Lookup("name").Value.String())
		//listStackResources()
		//sess := session.Must(session.NewSessionWithOptions(session.Options{
		//	SharedConfigState: session.SharedConfigEnable,
		//}))
		//// Create CloudFormation client
		//svc := cloudformation.New(sess)
		//parentStackResources, errParentStackResources := svc.DescribeStackResources(&cloudformation.DescribeStackResourcesInput{
		//	StackName: aws.String(stackName),
		//})
		//if errParentStackResources != nil {
		//	fmt.Printf("Stack with id %s does not exist\n", stackName)
		//	fmt.Println(errParentStackResources.Error())
		//	os.Exit(1)
		//}
		////Loop Parent parentStackResource resources
		//for _, parentStackResource := range parentStackResources.StackResources {
		//	fmt.Println("Stack Resource physical Id: " + *parentStackResource.PhysicalResourceId)
		//	fmt.Println("Stack Resource logical Id: " + *parentStackResource.LogicalResourceId)
		//	fmt.Println("Stack Name: " + *parentStackResource.StackName)
		//	fmt.Println(parentStackResource.String())
		//	svc.DescribeStacks(&cloudformation.DescribeStacksInput{
		//
		//	})
			// Fetch the resources for the child
			//childStackResources, errChildStackResources := svc.ListStackResources(&cloudformation.ListStackResourcesInput{
			//	StackName: aws.String(*parentStackResource.StackName),
			//})
			//if errChildStackResources != nil {
			//	fmt.Printf("Child Stack with id %s does not exist\n", aws.String(*parentStackResource.PhysicalResourceId))
			//	fmt.Println(errChildStackResources.Error())
			//	os.Exit(1)
			//}
			//for _, childstack := range childStackResources.StackResourceSummaries {
			//	fmt.Println("Child Stack Resource physical Id: " + *childstack.PhysicalResourceId)
			//	fmt.Println("Child Stack Resource logical Id: " + *childstack.LogicalResourceId)
			//	//fmt.Println("Child stack Name: " + *childstack.)
			//	fmt.Println(childstack.String())
			//}

			//fmt.Println(*parentStackResource.ResourceType + ", Status: " + *parentStackResource.ResourceStatus + ", Id: "+ *parentStackResource.PhysicalResourceId)
		//}
		// We skip DELETE_COMPLETE:
		//var filter = []*string{aws.String("CREATE_IN_PROGRESS"), aws.String("CREATE_FAILED"), aws.String("CREATE_COMPLETE"), aws.String("ROLLBACK_IN_PROGRESS"), aws.String("ROLLBACK_FAILED"), aws.String("ROLLBACK_COMPLETE"), aws.String("DELETE_IN_PROGRESS"), aws.String("DELETE_FAILED"), aws.String("UPDATE_IN_PROGRESS"), aws.String("UPDATE_COMPLETE_CLEANUP_IN_PROGRESS"), aws.String("UPDATE_COMPLETE"), aws.String("UPDATE_ROLLBACK_IN_PROGRESS"), aws.String("UPDATE_ROLLBACK_FAILED"), aws.String("UPDATE_ROLLBACK_COMPLETE_CLEANUP_IN_PROGRESS"), aws.String("UPDATE_ROLLBACK_COMPLETE"), aws.String("REVIEW_IN_PROGRESS")}
		//input := &cloudformation.ListStacksInput{StackStatusFilter: filter}
		//
		//parentStackResources, errParentStackResources := svc.ListStacks(input)
		//if errParentStackResources != nil {
		//	fmt.Println("Got error listing stacks:")
		//	fmt.Println(errParentStackResources.Error())
		//	os.Exit(1)
		//}

		//for _, parentStackResource := range parentStackResources.StackSummaries {
		//	fmt.Println(*parentStackResource.StackName + ", Status: " + *parentStackResource.StackStatus)
		//}
	//},
}

func init() {
	RootCmd.AddCommand(stackCmd)

	// Here you will define your flags and configuration settings.
	stackName = stackCmd.PersistentFlags().StringP("name","n","", "AWS Stack name")
	stackCmd.MarkPersistentFlagRequired("name")
	viper.BindEnv("name", stackCmd.PersistentFlags().Lookup("name").Value.String())

}

func listStackResources(cmd *cobra.Command, args []string) {
	unparsedResources := helpers.GetNestedCloudFormationResources(stackName)
	resources := make([]cfnResource, len(unparsedResources))

	c := make(chan cfnResource)
	for _, unparsedResource := range unparsedResources {
		go func(resource *cloudformation.StackResource) {
			resourceStruct := cfnResource{
				ResourceID:   aws.StringValue(resource.PhysicalResourceId),
				Type:         aws.StringValue(resource.ResourceType),
				Stack:        aws.StringValue(resource.StackName),
				Status:       aws.StringValue(resource.ResourceStatus),
				LogicalName:  aws.StringValue(resource.LogicalResourceId),
				ResourceName: aws.StringValue(resource.PhysicalResourceId),
			}
			// Override the resource name when there is a better name available
			switch resourceStruct.Type {
			case "AWS::EC2::Instance":
				resourceStruct.ResourceName = helpers.GetEc2Name(resource.PhysicalResourceId)
			case "AWS::RDS::DBInstance":
				resourceStruct.ResourceName = helpers.GetRDSName(resource.PhysicalResourceId)

			}
			c <- resourceStruct
		}(unparsedResource)
	}
	for i := 0; i < len(unparsedResources); i++ {
		resources[i] = <-c
	}
	keys := []string{"ResourceID", "Type", "Stack", "Name"}
	if *settings.Verbose {
		keys = append(keys, "Status")
		keys = append(keys, "LogicalName")
	}
	output := helpers.OutputArray{Keys: keys}
	for _, resource := range resources {
		content := make(map[string]string)
		content["ResourceID"] = resource.ResourceID
		content["Type"] = resource.Type
		content["Stack"] = resource.Stack
		content["Name"] = resource.ResourceName
		if *settings.Verbose {
			content["Status"] = resource.Status
			content["LogicalName"] = resource.LogicalName
		}
		holder := helpers.OutputHolder{Contents: content}
		output.AddHolder(holder)
	}
	output.Write(*settings)
}
