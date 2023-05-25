#!/bin/sh

queueName="${OKTETO_NAMESPACE}-oktacoshop"
createOutput=$(aws sqs create-queue --queue-name "$queueName" --tags owner="${OKTETO_NAME}")
exitCode=$?

if [ $exitCode -ne 0 ] && [ $exitCode -ne 254 ]; then
  echo "Failed to create SQS: exit code $exitCode"
  exit $exitCode
fi


if [ $exitCode -eq 0 ]; then
  queueUrl=$(printf %s "$createOutput" | jq '.["QueueUrl"]')
  echo "SQS queue ${queueUrl} created successfully"
fi

if [ $exitCode -eq 254 ]; then
  echo "SQS queue ${queueName} already exists"
  output=$(aws sqs get-queue-url --queue-name "$queueName" --output=json)
  queueUrl=$(echo "$output" | jq '.["QueueUrl"]')
fi

encoded=$(printf %s "$queue" | jq -sRr @uri)
dashboard="https://${AWS_REGION}.console.aws.amazon.com/sqs/v2/home?region=${AWS_REGION}#/queues/${encoded}"

{
  echo "OKTETO_EXTERNAL_SQS_ENDPOINTS_QUEUE_URL=$dashboard"
  echo "SQS_QUEUE_URL=$queueUrl"
  echo "SQS_QUEUE_NAME=$queueName"
} >> "$OKTETO_ENV"

echo "SQS queue configuration generated successfully"
exit 0


