createOutput=$(aws sqs create-queue --queue-name ${OKTETO_NAMESPACE}-oktacoshop --tags owner=${OKTETO_NAME})
exitCode=$?

if [ $exitCode -ne 0 && $exitCode -ne 254 ]; then
  echo "Failed to create SQS: exit code $exitCode"
  exit $exitCode
fi


if [ $exitCode -eq 0 ]; then
  echo "SQS queue created successfully"
  queue=$(echo $createOutput | jq '.["QueueUrl"]')
fi

if [ $exitCode -eq 254 ]; then
  echo "SQS queue already exists"
  output=$(aws sqs get-queue-url --queue-name ${OKTETO_NAMESPACE}-oktacoshop --output=json)
  queue=$(echo $output | jq '.["QueueUrl"]')
fi

encoded=$(printf %s "$queue" | jq -sRr @uri)
dashboard="https://${AWS_REGION}.console.aws.amazon.com/sqs/v2/home?region=${AWS_REGION}#/queues/${encoded}"


echo "OKTETO_EXTERNAL_SQS_ENDPOINTS_QUEUE_URL=${dashboard}" >> $OKTETO_ENV
echo "SQS_QUEUE_URL=${queue}" >> $OKTETO_ENV
echo "SQS_QUEUE_NAME=${OKTETO_NAMESPACE}-oktacoshop" >> $OKTETO_ENV
echo "SQS queue configuration generated successfully"
exit 0


