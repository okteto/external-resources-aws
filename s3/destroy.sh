aws s3api delete-bucket --bucket ${OKTETO_NAMESPACE}-oktacoshop > /dev/null
exitCode=$?
if [ $exitCode -eq 254 ]; then
  exit 0
fi

exit $exitCode