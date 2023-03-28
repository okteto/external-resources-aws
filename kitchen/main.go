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

var (
	//check origin will check the cross region source (note : please not using in production)
	upgrader                              = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	pendingOrders map[string]PendingOrder = map[string]PendingOrder{}
)

type FoodOrder struct {
	Items []string `json:"items"`
}

type FoodReady struct {
	OrderID string `json:"orderId"`
	Item    string `json:"item"`
}

type PendingOrder struct {
	OrderID string
	Items   []PendingOrderItem
}

type PendingOrderItem struct {
	Name  string
	Ready bool
}

func AddPendingOrder(orderId string, f FoodOrder) {
	p := PendingOrder{
		Items: make([]PendingOrderItem, len(f.Items)),
	}

	for i := range f.Items {
		p.Items[i] = PendingOrderItem{
			Name:  f.Items[i],
			Ready: false,
		}
	}

	pendingOrders[orderId] = p
}

func MarkItemReady(f FoodReady) {
	if k, ok := pendingOrders[f.OrderID]; ok {
		for i := range k.Items {
			if pendingOrders[f.OrderID].Items[i].Name == f.Item {
				pendingOrders[f.OrderID].Items[i].Ready = true
				fmt.Printf("Item '%s' from Order '%s' is ready 🍽️ \n", f.Item, f.OrderID)

			}
		}

		if k.IsReady() {
			fmt.Printf("Order '%s' is ready 🛎️!\n", f.OrderID)
		}
	}
}

func (p *PendingOrder) IsReady() bool {
	for _, i := range p.Items {
		if !i.Ready {
			return false
		}
	}

	return true
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

		MarkItemReady(ready)
		c.Status(200)
	})

	r.GET("/ws", func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer ws.Close()
		fmt.Printf("connected to the kitchen %s", *queueURL)
		fmt.Println()
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

			if len(msgResult.Messages) == 0 {
				continue
			}

			fmt.Printf("received %d messages from the queue", len(msgResult.Messages))
			fmt.Println()

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

				AddPendingOrder(*m.MessageId, order)

				fmt.Printf("sending order with %d items to the kitchen", len(order.Items))
				fmt.Println()

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
