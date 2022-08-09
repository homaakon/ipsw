/*
Copyright © 2022 blacktop

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/blacktop/ipsw/internal/utils"
	"github.com/blacktop/ipsw/pkg/usb/syslog"
	"github.com/caarlos0/ctrlc"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	idevCmd.AddCommand(iDevSyslogCmd)

	iDevSyslogCmd.Flags().StringP("uuid", "u", "", "Device UUID to connect")
	iDevSyslogCmd.Flags().DurationP("timeout", "t", time.Second*10, "Log timeout")
	iDevSyslogCmd.Flags().Bool("color", false, "Colorize output")
}

// iDevSyslogCmd represents the syslog command
var iDevSyslogCmd = &cobra.Command{
	Use:           "syslog",
	Short:         "Dump screenshot as a PNG",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		if Verbose {
			log.SetLevel(log.DebugLevel)
		}

		uuid, _ := cmd.Flags().GetString("uuid")
		timeout, _ := cmd.Flags().GetDuration("timeout")
		forceColor, _ := cmd.Flags().GetBool("color")

		color.NoColor = !forceColor

		if len(uuid) == 0 {
			dev, err := utils.PickDevice()
			if err != nil {
				return fmt.Errorf("failed to pick USB connected devices: %w", err)
			}
			uuid = dev.UniqueDeviceID
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := ctrlc.Default.Run(ctx, func() error {
			r, err := syslog.Syslog(uuid)
			if err != nil {
				return err
			}
			defer r.Close()
			_, err = io.Copy(os.Stdout, r)
			return err
		}); err != nil {
			return err
		}

		return nil
	},
}
