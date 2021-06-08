package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ExecuteCommand(dir string, name string, subname string, args ...string) (string, error) {
	args = append([]string{subname}, args...)
	Info.Printf("%s %s", name, args)
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		Error.Printf("CombinedOutput %s, failed with %s\n", string(bytes), err)
		os.Exit(1)
		//	fmt.Printf("CombinedOutput %s, failed with %s\n", string(bytes), err)
	}

	return string(bytes), err
}

func ExecuteError(cmd *cobra.Command, args []string, err error) {
	fmt.Fprintf(os.Stderr, "execute %s args:%v error:%v\n", cmd.Name(), args, err)
	os.Exit(1)
}

func PrintConfig() {
	fmt.Println(viper.Get("database.config"))
	fmt.Println(viper.Get("database.report_dir"))
	fmt.Println(viper.Get("database.work_dir"))
}

func GetDbConfig() (string, string, string) {
	config_dir := viper.Get("database.config")
	report_dir := viper.Get("database.report_dir")
	work_dir := viper.Get("database.work_dir")
	if report_dir == nil {
		fmt.Println("Error, report_dir not set")
		os.Exit(1)
	}
	return config_dir.(string), report_dir.(string), work_dir.(string)
}
