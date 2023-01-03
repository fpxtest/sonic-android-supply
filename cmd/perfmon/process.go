/*
 *   sonic-android-supply  Supply of ADB.
 *   Copyright (C) 2022  SonicCloudOrg
 *
 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU Affero General Public License as published
 *   by the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU Affero General Public License for more details.
 *
 *   You should have received a copy of the GNU Affero General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
package perfmon

import (
	"fmt"
	"github.com/SonicCloudOrg/sonic-android-supply/src/perfmonUtil"
	"github.com/SonicCloudOrg/sonic-android-supply/src/util"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"time"
)

var processPerfmonCmd = &cobra.Command{
	Use:   "process",
	Short: "Get app or pid performance",
	Long:  "Get app or pid performance",
	RunE: func(cmd *cobra.Command, args []string) error {
		if pid == "" && appName == "" {
			return fmt.Errorf("pid or app-name is require")
		}
		var err error
		device := util.GetDevice(serial)
		if pid == "" {
			pid, err = perfmonUtil.GetPidOnPackageName(device, appName)
			if err != nil {
				return err
			}
			if pid == "" {
				return fmt.Errorf("not find app corresponding pid")
			}
		}
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, os.Kill)
		timer := time.Tick(time.Duration(interval * int(time.Second)))
		done := false
		for !done {
			select {
			case <-sig:
				done = true
				fmt.Println()
			case <-timer:
				if processInfo, err := perfmonUtil.GetProcessInfo(device, pid, 1); err != nil {
					log.Panic(err)
				} else {
					if appName != "" {
						processInfo.Name = appName
					}
					data := util.ResultData(processInfo)
					fmt.Println(util.Format(data, isFormat, isJson))
				}
			}
		}
		return nil
	},
}

var appName string
var pid string

func initProcessPerfmon() {
	perfmonRootCMD.AddCommand(processPerfmonCmd)
	processPerfmonCmd.Flags().StringVarP(&appName, "app-name", "n", "", "applicationName")
	processPerfmonCmd.Flags().StringVarP(&pid, "pid", "p", "", "process id")
	processPerfmonCmd.Flags().StringVarP(&serial, "serial", "s", "", "device serial")
	processPerfmonCmd.Flags().IntVarP(&interval, "interval", "i", 1, "data refresh time")
	processPerfmonCmd.Flags().BoolVarP(&isJson, "json", "j", false, "convert to JSON string")
	processPerfmonCmd.Flags().BoolVarP(&isFormat, "format", "f", false, "convert to JSON string and format")
}
