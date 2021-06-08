/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// databasesCmd represents the databases command
var databasesCmd = &cobra.Command{
	Use:   "databases",
	Short: "databases support On-Site Inspection",
	Long:  `databases support On-Site Inspection`,
	// not support databasesCmd.Flags().StringP
	//	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		parameter, _ := cmd.Flags().GetString("database_tpye")
		DbExecuteCommand(parameter)
	},
}

func DbExecuteCommand(parameter string) {
	if parameter == "mysql" {
		config_dir, report_dir, work_dir := GetDbConfig()
		format_config_dir := strings.Replace("--defaults-file=/tmp/my.cnf", "/tmp/my.cnf", config_dir, 1)
		format_report_dir := strings.Replace("source /tmp/mysqlreport.sql", "/tmp/mysqlreport.sql", report_dir, 1)
		output, _ := ExecuteCommand(work_dir, "mysql", format_config_dir, "-e", format_report_dir)
		//if err != nil {
		//	Error.Printf("execute %s args:%v error: %v\n", cmd.Name(), args, err)
		//	ExecuteError(cmd, args, err)

		//}
		Info.Printf(output)
		//	fmt.Fprint(os.Stdout, output)

	} else if parameter == "oracle" {
		_, report_dir, work_dir := GetDbConfig()
		output, _ := ExecuteCommand(work_dir, "sqlplus", "/ as sysdba", report_dir)
		Info.Printf(output)
	} else {
		Error.Println("not support")
	}
}

func init() {
	// add local sub cmd && flag
	databasesCmd.Flags().StringP("database_tpye", "d", "mysql", "mysql or oracle")
	databasesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	viper.BindPFlag("database_tpye", databasesCmd.Flags().Lookup("database_tpye"))
	rootCmd.AddCommand(databasesCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// databasesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
