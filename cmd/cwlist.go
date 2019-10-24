// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
	"log"
	"regexp"
)

// listCmd represents the list command
var cwListCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")
		apply, _ := cmd.Flags().GetBool("apply")
		retention, _ := cmd.Flags().GetInt64("retention")
		namePattern, _ := cmd.Flags().GetString("name")

		listLogGroups(retention, apply, namePattern)
	},
}

func listLogGroups(retention int64, apply bool, namePattern string) {
	log.Printf("Retention days: %v", retention)
	log.Printf("Apply: %v", apply)
	//Create a session
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	region := *sess.Config.Region
	log.Printf("Running on region: %v", region)
	//Create a service
	svc := cloudwatchlogs.New(sess)
	// Create an input
	input := &cloudwatchlogs.DescribeLogGroupsInput{}

	result, err := svc.DescribeLogGroups(input)
	if err != nil {
		log.Println(err)
		return
	}

	cloudwatchGroups := result.LogGroups

	for result.NextToken != nil {
		input = &cloudwatchlogs.DescribeLogGroupsInput{
			NextToken: result.NextToken,
		}
		result, err = svc.DescribeLogGroups(input)
		if err != nil {
			log.Println(err)
			return
		}
		for _, group := range result.LogGroups {
			cloudwatchGroups = append(cloudwatchGroups, group)
		}
	}

	// Calculate the region total size for cost
	var totalLogByteSize int64
	noRetentionSet := false

	for _, group := range cloudwatchGroups {
		//LogGroupName := "/aws/lambda/cilogging20/error-logs"
		match, _:= regexp.MatchString(namePattern, *group.LogGroupName)
		if match {
			totalLogByteSize = totalLogByteSize + *group.StoredBytes
			noRetentionSet = true
			if apply == true {
				// set input filter
				input := &cloudwatchlogs.PutRetentionPolicyInput{
					LogGroupName:    aws.String(*group.LogGroupName),
					RetentionInDays: aws.Int64(retention),
				}
				// put retention policy
				_, err := svc.PutRetentionPolicy(input)
				if err != nil {
					log.Println(err)
					return
				}
				log.Printf("Retention policy for %s was set to %v", *group.LogGroupName, retention)
			} else {
				log.Printf("Group %s retention policy would be set to %d (size is %v Bytes), --yes to apply", *group.LogGroupName, retention, *group.StoredBytes)
			}
		}
	}
	if noRetentionSet == true {
		log.Printf("Region %s total log size:", region)
		log.Printf("Total log size in with no retention policy is: %v bytes", totalLogByteSize)
	}

	fmt.Println("=====================================================================================================")
}

func init() {
	cwCmd.AddCommand(cwListCmd)
	//cwCmd.Flags().Int64("retention", 1, "Setting the retention policy in days")
	//cwCmd.Flags().Bool("apply", false, "Apply changes.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	cwListCmd.Flags().BoolP("apply", "y", false, "Update loggroup")
	cwListCmd.Flags().Int64P("retention","r", 1, "Default retention days.")
	cwListCmd.Flags().String("name", "/aws/lambda/cilogging20-", "Log group name prefix pattern." )
}
