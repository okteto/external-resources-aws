const express = require('express')
const { SQSClient, SendMessageCommand } = require("@aws-sdk/client-sqs");
const sqsClient = new SQSClient({region: process.env.REGION});
const author = process.env.AUTHOR;
const queue = process.env.QUEUE;

const app = express()
app.use(express.json());
app.use(express.static('public'))

app.get('/healthz', function (req, res) {
  res.sendStatus(200);
})

app.post('/order', function (req, res) {
  console.log(req.body);
  var params = {
   MessageBody: JSON.stringify(req.body.order),
   QueueUrl: queue
 };

 sqsClient.send(new SendMessageCommand(params))
  .then(data => {
    console.log(data);
    res.sendStatus(201);
  })
  .catch(error => {
    console.error(error);
    res.sendStatus(500);
  });
})


console.log(`hi, my name is ${author}`);
app.listen(3000)