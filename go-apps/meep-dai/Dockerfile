# Copyright (c) 2022  The AdvantEDGE Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


FROM debian:9.6-slim

COPY ./meep-dai /meep-dai
COPY ./api /api
COPY ./user-api /user-api
COPY ./data /
#COPY ./onboardedapp/ /onboardedapp

VOLUME /onboardedapp-vol

RUN chmod +x /entrypoint.sh

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends procps

EXPOSE 31120-31129/tcp

ENTRYPOINT ["/entrypoint.sh"]
