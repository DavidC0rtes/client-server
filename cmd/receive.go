/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	client_tcp "github.com/DavidC0rtes/client-server/client"

	"github.com/spf13/cobra"
)

var channelReceive int

// receiveCmd represents the receive command
var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "Starts the client on receive mode.",
	Long: `Clients can subscribe channels by starting in receive mode. Every byte sent to the specified channel
will be delivered to said listening client.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("receive called")
		client_tcp.Subscribe(channelReceive)
	},
}

func init() {
	clientCmd.AddCommand(receiveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// receiveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// receiveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	receiveCmd.Flags().IntVarP(&channelReceive, "channel", "c", 0, "Channel to listen to.")
	receiveCmd.MarkFlagRequired("channel")
}
