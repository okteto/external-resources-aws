#!/bin/sh

bucketName="${OKTETO_NAMESPACE}-oktacoshop"
aws s3api delete-bucket --bucket "$bucketName"
exitCode=$?
if [ $exitCode -eq 254 ] || [ $exitCode -eq 0 ]; then
  echo "S3 bucket deleted successfully"
  exit 0
fi

exit $exitCode