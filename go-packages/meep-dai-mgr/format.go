/*
 * Copyright (c) 2022  The AdvantEDGE Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package meepdaimgr

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
