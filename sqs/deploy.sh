createOutput=$(aws sqs create-queue --queue-name ${OKTETO_NAMESPACE}-oktacoshop --tags owner=${OKTETO_NAME})
exitCode=$?

if [ $exitCode -eq 0 ]; then
  queue=$(echo $createOutput | jq '.["QueueUrl"]')
  echo "OKTETO_EXTERNAL_SQS_ENDPOINTS_QUEUE_URL=${queue}" >> $OKTETO_ENV
  echo "SQS_QUEUE_NAME=${OKTETO_NAMESPACE}-oktacoshop" >> $OKTETO_ENV
  exit 0
fi

if [ $exitCode -eq 254 ]; then
  output=$(aws sqs get-queue-url --queue-name ${OKTETO_NAMESPACE}-oktacoshop --output=json)
  queue=$(echo $output | jq '.["QueueUrl"]')
  echo "OKTETO_EXTERNAL_SQS_ENDPOINTS_QUEUE_URL=${queue}" >> $OKTETO_ENV
  echo "SQS_QUEUE_NAME=${OKTETO_NAMESPACE}-oktacoshop" >> $OKTETO_ENV
  exit 0
fi

exit $exitCode