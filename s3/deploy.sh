bucketName=${OKTETO_NAMESPACE}-oktacoshop
createOutput=$(aws s3 mb s3://$bucketName)
exitCode=$?

if [ $exitCode -ne 0 && $exitCode -ne 254 ]; then
  echo "Failed to create S3: exit code $exitCode"
  exit $exitCode
fi


if [ $exitCode -eq 0 ]; then
  echo "S3 bucket created successfully"
fi

if [ $exitCode -eq 254 ]; then
  echo "S3 bucket already exists"
fi

dashboard="https://s3.console.aws.amazon.com/s3/buckets/$bucketName"


echo "OKTETO_EXTERNAL_S3_ENDPOINTS_BUCKET_URL=${dashboard}" >> $OKTETO_ENV
echo "S3_BUCKET_NAME=$bucketName" >> $OKTETO_ENV
echo "S3 bucket configuration generated successfully"
exit 0


