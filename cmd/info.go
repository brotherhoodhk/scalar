/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

type info struct {
	XMLName xml.Name   `xml:"scalar"`
	Version string     `xml:"version,attr"`
	Basic   basic_info `xml:"basic_info"`
}
type basic_info struct {
	XMLName xml.Name `xml:"basic_info"`
	Site    string   `xml:"offical_site"`
}

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "scalar based on gocache",
	Long: `scalar is a cache system,based on gocache. it can clear the useless data when data volume reaches the threshold.
And scalar has many zones,you can put your data in the zone that you want`,
	Run: func(cmd *cobra.Command, args []string) {
		fe, err := ioutil.ReadFile(os.Getenv("SCALAR_HOME") + "/conf/info.xml")
		infomation := new(info)
		if err == nil {
			err = xml.Unmarshal(fe, infomation)
		}
		if err == nil {
			fmt.Printf("version:%v\noffical site:%v\n", infomation.Version, infomation.Basic.Site)
		} else {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
