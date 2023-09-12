/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"scal/basic"

	driver_tools "github.com/oswaldoooo/gocache-driver/basics"
	"github.com/spf13/cobra"
)

type siteconfig struct {
	XMLName   xml.Name    `xml:"scalar"`
	Cacheinfo gocacheinfo `xml:"gocache"`
}
type gocacheinfo struct {
	XMLName    xml.Name `xml:"gocache"`
	Host       string   `xml:"hostadd"`
	Port       int      `xml:"port"`
	Default_DB string   `xml:"default_db"`
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "scalar based on gocache",
	Long: `scalar is a cache system,based on gocache. it can clear the useless data when data volume reaches the threshold.
And scalar has many zones,you can put your data in the zone that you want`,
	Run: func(cmd *cobra.Command, args []string) {
		buff, err := ioutil.ReadFile(os.Getenv("SCALAR_HOME") + "/conf/site.xml")
		if err == nil {
			siteconf := new(siteconfig)
			err = xml.Unmarshal(buff, siteconf)
			if err == nil && len(siteconf.Cacheinfo.Host) > 0 && len(siteconf.Cacheinfo.Default_DB) > 0 && siteconf.Cacheinfo.Port > 0 {
				dbcon := driver_tools.New(siteconf.Cacheinfo.Host, siteconf.Cacheinfo.Port, "", siteconf.Cacheinfo.Default_DB)
				err = dbcon.Connect()
				if err == nil {
					defer dbcon.Close()
					err = dbcon.CreateDB()
					if err == nil {
						err = dbcon.SetKey("version", basic.VERSION)
						if err == nil {
							err = dbcon.Save()
						}
					}
				}
			}
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		} else {
			fmt.Println("init success")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
