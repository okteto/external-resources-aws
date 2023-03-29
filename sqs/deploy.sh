createOutput=$(aws sqs create-queue --queue-name ${OKTETO_NAMESPACE}-oktacoshop --tags owner=${OKTETO_NAME})
exitCode=$?

if [ $exitCode -eq 0 ]; then
  echo "SQS queue created successfully"
  queue=$(echo $createOutput | jq '.["QueueUrl"]')
  echo "OKTETO_EXTERNAL_SQS_ENDPOINTS_QUEUE_URL=${queue}" >> $OKTETO_ENV
  echo "SQS_QUEUE_NAME=${OKTETO_NAMESPACE}-oktacoshop" >> $OKTETO_ENV
  echo "SQS queue configuration generated successfully"
  exit 0
fi

if [ $exitCode -eq 254 ]; then
  echo "SQS queue already exists"
  output=$(aws sqs get-queue-url --queue-name ${OKTETO_NAMESPACE}-oktacoshop --output=json)
  queue=$(echo $output | jq '.["QueueUrl"]')
  echo "OKTETO_EXTERNAL_SQS_ENDPOINTS_QUEUE_URL=${queue}" >> $OKTETO_ENV
  echo "SQS_QUEUE_NAME=${OKTETO_NAMESPACE}-oktacoshop" >> $OKTETO_ENV
  echo "SQS queue configuration generated successfully"
  exit 0
fi

echo "Failed to create SQS: exit code $exitCode"
exit $exitCode