/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"

	server_tcp "github.com/DavidC0rtes/client-server/server-tcp"
)

var numChannels int
var maxFilesize int64

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Initialize the server",
	Long:  `tells the application to start in server mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		server_tcp.Run(numChannels, maxFilesize)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serverCmd.Flags().IntVarP(&numChannels, "channels", "c", 3, "Number of channels to create.")
	serverCmd.Flags().Int64VarP(&maxFilesize, "max", "m", 4096, "Maximum supported filesize (B).")
}
