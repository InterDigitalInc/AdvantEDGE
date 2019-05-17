/*
 * Copyright (c) 2019
 * InterDigital Communications, Inc.
 * All rights reserved.
 *
 * The information provided herein is the proprietary and confidential
 * information of InterDigital Communications, Inc.
 */

package utils

import (
	"time"

	"github.com/spf13/cobra"
)

// FormatResult consistent formatting for results to be displayed
func FormatResult(result string, elapsed time.Duration, cobraCmd *cobra.Command) string {
	ret := result

	t, _ := cobraCmd.Flags().GetBool("time")
	if t {
		ret += (" [" + elapsed.Round(time.Millisecond).String() + "]")
	}

	return ret
}
