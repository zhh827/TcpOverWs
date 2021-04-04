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
	"log"
	"tcpoverws/service"

	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client  -w wsurl，exp: tcpoverws client  -w ws://a.com/ -t 192.168.1.1:22",
	Short: "tcp server -> websocket client",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("[INFO] Run in client mode, tcp listen on: %s \n", tcplisten)
		log.Printf("[INFO] Connect to %s\n", wstarget)
		service.Client(wstarget, tcplisten, tcptarget, passwd)
	},
}

var tcplisten string
var wstarget string
var tcptarget string

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.Flags().StringVarP(&tcplisten, "listen", "l", "0.0.0.0:5001", "locla listen port")
	clientCmd.Flags().StringVarP(&passwd, "passwd", "p", "syswin@123", "login password")
	clientCmd.Flags().StringVarP(&wstarget, "ws", "w", "", "connect websocket server(required)")
	clientCmd.Flags().StringVarP(&tcptarget, "target", "t", "", "taget IP:PORT(required)")
	clientCmd.MarkFlagRequired("ws")
	clientCmd.MarkFlagRequired("target")
}
