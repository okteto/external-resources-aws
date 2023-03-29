package main

import (
	"bytes"
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
	OrderID string             `json:"orderId"`
	Items   []PendingOrderItem `json:"items"`
}

type PendingOrderItem struct {
	Name  string `json:"name"`
	Ready bool   `json:"ready"`
}

func CreatePendingOrder(orderID string, f FoodOrder) PendingOrder {
	p := PendingOrder{
		OrderID: orderID,
		Items:   make([]PendingOrderItem, len(f.Items)),
	}

	for i := range f.Items {
		p.Items[i] = PendingOrderItem{
			Name:  f.Items[i],
			Ready: false,
		}
	}

	pendingOrders[orderID] = p
	return p
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
			k.OrderCheck()
		}
	}
}

func (p *PendingOrder) OrderCheck() {
	checkServiceUrl := os.Getenv("CHECK")
	buff := new(bytes.Buffer)
	json.NewEncoder(buff).Encode(p)

	_, err := http.Post(checkServiceUrl, "application/json", buff)
	if err != nil {
		fmt.Printf("failed to order check: %s", err)
		fmt.Println()
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
				WaitTimeSeconds:     aws.Int64(3),
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

				var order FoodOrder
				if err := json.Unmarshal([]byte(*m.Body), &order); err != nil {
					fmt.Printf("failed to unmarshall the message: %s", err)
					fmt.Println()
					break
				}

				p := CreatePendingOrder(*m.MessageId, order)

				fmt.Printf("sending order %s with %d items to the kitchen", p.OrderID, len(p.Items))
				fmt.Println()

				//Response message to client
				err = ws.WriteJSON(p)
				if err != nil {
					fmt.Printf("failed to send the message to the socket: %s", err)
					fmt.Println()
					break
				}

				svc.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:      queueURL,
					ReceiptHandle: m.ReceiptHandle,
				})

				fmt.Printf("completed message %s ", *m.ReceiptHandle)
				fmt.Println()

			}
		}
	})

	r.StaticFS("/public", http.Dir("public"))
	r.StaticFile("/", "./public/index.html")

	fmt.Println("ready to cook some grub 🔪")
	r.Run()
}
