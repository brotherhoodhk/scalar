/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"scal/basic"

	"github.com/spf13/cobra"
)

var reinitall = "all"

// reinitCmd represents the reinit command
var reinitCmd = &cobra.Command{
	Use:   "reinit",
	Short: "scalar based on gocache",
	Long: `scalar is a cache system,based on gocache. it can clear the useless data when data volume reaches the threshold.
And scalar has many zones,you can put your data in the zone that you want`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(reinitall) == 0 || reinitall == " " {
			reinitall = "all"
		}
		switch reinitall {
		case "all":
			errone := reinitFilemap()
			errtwo := reinitDblink()
			if errone != nil {
				fmt.Println(errone)
			}
			if errtwo != nil {
				fmt.Println(errtwo)
			}
		case "filemap":
			if err := reinitFilemap(); err != nil {
				fmt.Println(err)
			}
		case "dblink":
			if err := reinitDblink(); err != nil {
				fmt.Println(err)
			}
		default:
			fmt.Println("dont support reinit", reinitall)
		}
	},
}

func init() {
	reinitCmd.Flags().StringVarP(&reinitall, "filemap", "m", "", "--filemap")
	rootCmd.AddCommand(reinitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reinitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reinitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func reinitFilemap() (err error) {
	mappath := basic.ROOTPATH + "/conf/filemap"
	fe, err := os.OpenFile(mappath, os.O_WRONLY|os.O_TRUNC, 0600)
	if err == nil {
		_, err = fe.Write([]byte(""))
	}
	return
}
func reinitDblink() (err error) {
	linkpath := basic.ROOTPATH + "/dblink"
	err = os.RemoveAll(linkpath)
	os.Mkdir(linkpath, 0600)
	return
}
