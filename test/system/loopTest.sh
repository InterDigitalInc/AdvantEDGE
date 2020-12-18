#!/bin/bash

# This script executes go test X times
export MEEP_HOST_TEST_URL="http://10.3.16.150"

i="0"

while [ $i -lt 66 ]
do
output_prefix="result"
output_suffix=".txt"
output=${output_prefix}${i}${output_suffix}
go test -timeout 20m > ${output}
i=$[$i+1]
done
