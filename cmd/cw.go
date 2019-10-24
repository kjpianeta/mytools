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

	"github.com/spf13/cobra"
)

// cwCmd represents the cw command
var cwCmd = &cobra.Command{
	Use:   "cw",
	Short: "Cloudwatch commands",
	Long: `This command enables you to execute various Cloudwatch related commands. Please review available sub-commands`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cw called")
		fmt.Println("Loggroup name: " + *loggroupname)
	},
}

var loggroupname *string
func init() {
	RootCmd.AddCommand(cwCmd)
	loggroupname = 	cwCmd.PersistentFlags().StringP("loggroup", "l", "DefaultLogGroup", "The name of the log group")
}
