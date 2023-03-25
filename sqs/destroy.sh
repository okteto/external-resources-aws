output=$(aws sqs get-queue-url --queue-name ${OKTETO_NAMESPACE}-oktacoshop --output=json)
if [ $? -eq 254 ]; then
  exit 0
fi

queue=$(echo $output |  jq -r '.["QueueUrl"]')
aws sqs delete-queue --queue-url $queue > /dev/null
exitCode=$?
if [ $exitCode -eq 254 ]; then
  exit 0
fi

exit $exitCode