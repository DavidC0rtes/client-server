/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/DavidC0rtes/client-server/client-tcp"
	"github.com/spf13/cobra"
)

var channel int

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use: "send",
	Args: func(sendCmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Requires filename argument.")
		} else if len(args) > 1 {
			return errors.New("Too many arguments.")
		}
		return nil
	},
	Short: "Starts the client in sending mode.",
	Long:  `Allows the client to send files to the TCP server in the specified channel.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("send called")
		client_tcp.PrepareSend(args[0], channel)
	},
}

func init() {
	clientCmd.AddCommand(sendCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sendCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sendCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	sendCmd.Flags().IntVarP(&channel, "channel", "c", -1, "Specify channel id")
	sendCmd.MarkFlagRequired("channel")
}
