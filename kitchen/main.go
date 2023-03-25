package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type FoodOrder struct {
	Order []string `json:"order"`
}

type FoodReady struct {
	OrderID string `json:"orderId`
	Item    string `json:"item"`
}

var upgrader = websocket.Upgrader{
	//check origin will check the cross region source (note : please not using in production)
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {

	sess := session.Must(session.NewSessionWithOptions(session.Options{}))
	svc := sqs.New(sess)
	urlResult, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(os.Getenv("QUEUE")),
	})

	if err != nil {
		panic(err)
	}

	queueURL := urlResult.QueueUrl

	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.POST("/ready", func(c *gin.Context) {
		var ready FoodReady
		if err := c.BindJSON(&ready); err != nil {
			fmt.Println(err)
			c.AbortWithStatus(500)
		}

		fmt.Printf("Item '%s' from Order '%s' is ready", ready.Item, ready.OrderID)
		c.Status(200)
	})

	r.GET("/ws", func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer ws.Close()
		for {
			//Read Message from the SQS queue
			msgResult, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
				QueueUrl:            queueURL,
				MaxNumberOfMessages: aws.Int64(1),
				WaitTimeSeconds:     aws.Int64(5),
			})

			if err != nil {
				fmt.Println(err)
				break
			}

			for _, m := range msgResult.Messages {
				svc.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:      queueURL,
					ReceiptHandle: m.ReceiptHandle,
				})

				var order FoodOrder
				if err := json.Unmarshal([]byte(*m.Body), &order); err != nil {
					fmt.Println(err)
					break
				}

				//Response message to client
				err = ws.WriteJSON(order)
				if err != nil {
					fmt.Println(err)
					break
				}
			}

		}
	})

	r.StaticFS("/public", http.Dir("public"))
	r.StaticFile("/", "./public/index.html")

	fmt.Println("ready to cook some grub 🔪")
	r.Run()
}
