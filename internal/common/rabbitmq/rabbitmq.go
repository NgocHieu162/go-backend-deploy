package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-backend/internal/common/env"
	"log"
	"math/rand"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn *amqp.Connection
}

func NewRabbitMQ(env *env.Env) *RabbitMQ {
	conn, err := amqp.DialConfig(env.RabbitMQURL, amqp.Config{
		Properties: amqp.Table{
			"connection_name": "go-backend",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("[RABBIT_MQ] Connect to rabbitmq successfully")
	return &RabbitMQ{
		Conn: conn,
	}
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func queueDeclareMain(ch *amqp.Channel, queueName string) (amqp.Queue, error) {
	return ch.QueueDeclare(
		queueName, // name
		true,      // durability: có tồn tại queue hay không khi rabbitmq restart
		false,     // delete when unused, vì là queue chính nên không tự động xóa
		false,     // exclusive: độc quyền, có khóa với connection hiện tại hay không, vì là còn B xử lý => false
		false,     // noWait có đợi queue tạo thành công hay không
		nil,       // arguments
	)
}

func queueDeclareTemp(ch *amqp.Channel, queueName string) (amqp.Queue, error) {
	return ch.QueueDeclare(
		queueName, // name
		true,      // durability: có tồn tại queue hay không khi rabbitmq restart
		true,      // delete when unused, vì là queue tạm cho nên phải tự động xóa
		false,     // exclusive: độc quyền, có khóa với connection hiện tại hay không, vì là còn B xử lý => false
		false,     // noWait có đợi queue tạo thành công hay không
		nil,       // arguments
	)
}

func publishNotReply(ctx context.Context, ch *amqp.Channel, queueName string, corrId string, body []byte) error {
	return ch.PublishWithContext(
		ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrId,

			// với những message quan trọng thì không chỉ lưu trữ trong RAM
			// giúp lưu lại message khi rabbitmq restart
			// đảm bảo queue sẽ được giữ: durable: true
			DeliveryMode: amqp.Persistent,
			// ReplyTo:       queueMain.Name,
			Body: body,
		})
}

func publishWithReply(ctx context.Context, ch *amqp.Channel, queueName string, corrId string, replyTo string, body []byte) error {
	return ch.PublishWithContext(
		ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrId,

			// với những message quan trọng thì không chỉ lưu trữ trong RAM
			// giúp lưu lại message khi rabbitmq restart
			// đảm bảo queue sẽ được giữ: durable: true
			DeliveryMode: amqp.Persistent,
			ReplyTo:      replyTo,
			Body:         body,
		})
}

type Reply struct {
	ErrorString string
	Data        json.RawMessage
}

func (r *RabbitMQ) Send(ctx context.Context, queueName string, payload any) (err error) {
	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	ch, err := r.Conn.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}

	// tạo queue nếu chưa có, nếu đã có sẽ kiểm tra đúng các setting hay không
	queueMain, err := queueDeclareMain(ch, queueName)
	if err != nil {
		fmt.Println(err)
		return
	}

	corrId := randomString(32)
	err = publishNotReply(ctx, ch, queueMain.Name, corrId, body)
	if err != nil {
		fmt.Println(err)
		return
	}

	return nil
}

func (r *RabbitMQ) On(ctx context.Context, queueName string, handler func(context.Context, []byte) error) (err error) {
	ch, err := r.Conn.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}

	// tạo queue nếu chưa có, nếu đã có sẽ kiểm tra đúng các setting hay không
	queueMain, err := queueDeclareMain(ch, queueName)

	if err != nil {
		fmt.Println(err)
		return
	}

	msgs, err := ch.Consume(
		queueMain.Name, // queue
		"",             // consumer
		false,          // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		defer ch.Close()
		for d := range msgs {

			//handler
			err = handler(ctx, d.Body)
			if err != nil {
				d.Nack(false, false)
				continue
			}

			d.Ack(false)

		}
	}()

	return nil
}

func (r *RabbitMQ) Request(ctx context.Context, queueName string, payload any, result any) (err error) {
	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	ch, err := r.Conn.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ch.Close()

	//==============================================
	queueTemp, err := queueDeclareTemp(ch, "")
	if err != nil {
		fmt.Println(err)
		return
	}

	msgs, err := ch.Consume(
		queueTemp.Name, // queue
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	//==============================================
	// tạo queue nếu chưa có, nếu đã có sẽ kiểm tra đúng các setting hay không
	queueMain, err := queueDeclareMain(ch, queueName)
	if err != nil {
		fmt.Println(err)
		return
	}

	corrId := randomString(32)
	err = publishWithReply(ctx, ch, queueMain.Name, corrId, queueTemp.Name, body)
	if err != nil {
		fmt.Println(err)
		return
	}

	//==============================

	for {
		select {
		case <-ctx.Done():
			err = errors.New("Request timeout")
			fmt.Println(err)
			return

		case d := <-msgs:
			if d.CorrelationId != corrId {
				continue
			}

			var replyBody Reply

			err = json.Unmarshal(d.Body, &replyBody)
			if err != nil {
				fmt.Println(err)
				return
			}

			if replyBody.ErrorString != "" {
				err = errors.New(replyBody.ErrorString)
				fmt.Println(err)
				return
			}

			err = json.Unmarshal(replyBody.Data, result)
			if replyBody.ErrorString != "" {
				fmt.Println(err)
				return
			}

			return
		}
	}
}

func (r *RabbitMQ) OnReply(ctx context.Context, queueName string, handler func(context.Context, []byte) (any, error)) (err error) {
	ch, err := r.Conn.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}

	// tạo queue nếu chưa có, nếu đã có sẽ kiểm tra đúng các setting hay không
	queueMain, err := queueDeclareMain(ch, queueName)
	if err != nil {
		fmt.Println(err)
		return
	}

	msgs, err := ch.Consume(
		queueMain.Name, // queue
		"",             // consumer
		false,          // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		defer ch.Close()
		for d := range msgs {
			// tạo hàm trả lỗi cho các logic phía dưới
			replyErr := func(replyError error) {
				d.Nack(false, false)
				replyBody := Reply{
					ErrorString: replyError.Error(),
					Data:        nil,
				}
				body, err := json.Marshal(replyBody)
				if err != nil {
					fmt.Println(err)
					return
				}

				err = publishNotReply(ctx, ch, d.ReplyTo, d.CorrelationId, body)
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			//handler
			result, err := handler(ctx, d.Body)
			if err != nil {
				replyErr(err)
				continue
			}

			data, err := json.Marshal(result)
			if err != nil {
				replyErr(err)
				continue
			}

			replyBody := Reply{
				ErrorString: "",
				Data:        data,
			}

			body, err := json.Marshal(replyBody)
			if err != nil {
				replyErr(err)
				continue
			}

			err = publishNotReply(ctx, ch, d.ReplyTo, d.CorrelationId, body)
			if err != nil {
				d.Nack(false, false)
				continue
			}

			d.Ack(false)
		}
	}()

	return nil
}
