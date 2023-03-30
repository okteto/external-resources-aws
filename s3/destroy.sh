#!/bin/sh

bucketName="${OKTETO_NAMESPACE}-oktacoshop"
aws s3 rm "s3://$bucketName" --recursive

output=$(aws s3 rb "s3://$bucketName")
exitCode=$?
echo $exitCode
if [ $exitCode -eq 0 ]; then
  echo "S3 bucket deleted successfully"
  exit 0
fi

echo "$output" | grep -q "NoSuchBucket" 
if [ $? -eq 1 ]; then
    exit 0
fi

exit 1