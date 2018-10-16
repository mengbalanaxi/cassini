package msgqueue

import (
	"github.com/nats-io/go-nats"
	"github.com/QOSGroup/cassini/log"
	"time"
	"errors"
)

//type Producer interface {
//	Produce(nc *nats.Conn, msg []byte) error
//}

type NATSProducer struct {
	ServerUrls string //消息队列服务地址，多个用","分割  例如 "nats://192.168.168.195:4222，nats://192.168.168.195:4223"
	Subject    string //主题
}

func (n *NATSProducer) Connect() (nc *nats.Conn,err error){
	nc, err = nats.Connect(n.ServerUrls)
	if err != nil {
		log.Error("Can't connect: %v\n", err)
		return nil,err
	}
	return
}

func (n *NATSProducer) Produce(nc *nats.Conn, msg []byte) (err error) {
	if nc == nil {
		return errors.New("the nats.Conn is nil")
	}
	//reconnect to nats server
	i := nc.Status()
	if i != 1 { //TODO 把1变成 nc.Status 的 常量 “CONNECTED”
		if i != 2 {nc.Close()} //status==2 closed
		nc, err = n.Connect()
		if err != nil {
			return errors.New("the nats.Conn is not available")
		}
	}

	if e := nc.Publish(n.Subject, msg);e != nil {
		return errors.New("Published faild" )
	}
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Error(err)
	} else {
		log.Infof("Published [%s] : '%T'\n", n.Subject, msg)
	}
	return nil
}

//TODO
func (n *NATSProducer) ProduceWithReply(nc *nats.Conn,reply string, payload []byte) error {

	msg, err := nc.Request(n.Subject,payload, 100*time.Millisecond)
	if err != nil {
		if nc.LastError() != nil {
			log.Errorf("Error in Request: %v\n", nc.LastError())
		}
		log.Errorf("Error in Request: %v\n", err)
	}

	log.Infof("Published [%s] : '%s'\n", n.Subject, payload)
	log.Infof("Received [%v] : '%s'\n", msg.Subject, string(msg.Data))
	log.Infof("Reply [%v] : '%s'\n", msg.Subject, string(msg.Reply))
	return nil
}