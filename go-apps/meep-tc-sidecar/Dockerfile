# Copyright (c) 2019
# InterDigital Communications, Inc.
# All rights reserved.
#
# The information provided herein is the proprietary and confidential
# information of InterDigital Communications, Inc.

FROM debian:9.6-slim
RUN apt-get update && apt-get install -y iputils-ping iproute2 iptables conntrack
COPY ./meep-tc-sidecar /meep-tc-sidecar
ENTRYPOINT /meep-tc-sidecar
