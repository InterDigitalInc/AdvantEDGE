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

package helm

func GetReleasesName(sandboxName string) ([]Release, error) {
	return getReleasesName(sandboxName)
}

/*
* currently GetReleases is not used. Since it uses helm status and helmv3 doesn't show resources
* https://github.com/helm/helm/issues/5952
 */

func GetReleases(sandboxName string) ([]Release, error) {
	return getReleases(sandboxName)
}

func InstallCharts(charts []Chart, sandboxName string) error {
	return runTask(Install, charts, sandboxName)
}

func DeleteReleases(charts []Chart, sandboxName string) error {
	return runTask(Delete, charts, sandboxName)
}
