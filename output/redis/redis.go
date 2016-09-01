package outputredis

import (
	"errors"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"

	"logcool/utils"
	"logcool/utils/logevent"
)

const (
	ModuleName = "redis"
)

type OutputConfig struct {
	utils.OutputConfig
	Key               string `json:"key"`
	Host              string `json:"host"`
	Password          string `json:"password"`
	DataType          string `json:"data_type"`
	Timeout           int    `json:"timeout"`
	ReconnectInterval int    `json:"reconnect_interval"`

	conns  *redis.Pool
	evchan chan logevent.LogEvent
}

func InitHandler(confraw *utils.ConfigRaw) (retconf utils.TypeOutputConfig, err error) {
	conf := OutputConfig{
		OutputConfig: utils.OutputConfig{
			CommonConfig: utils.CommonConfig{
				Type: ModuleName,
			},
		},
		Key:               "logcool",
		DataType:          "list",
		Timeout:           5,
		ReconnectInterval: 1,

		evchan: make(chan logevent.LogEvent),
	}
	if err = utils.ReflectConfig(confraw, &conf); err != nil {
		return
	}
	conf.conns = initRedisPool(conf.Host, conf.Password)

	go conf.loop()

	retconf = &conf
	return
}

func (self *OutputConfig) Event(event logevent.LogEvent) (err error) {
	self.evchan <- event
	return
}

func (self *OutputConfig) loop() (err error) {
	for {
		event := <-self.evchan
		self.sendEvent(event)
	}

	return
}

func (self *OutputConfig) sendEvent(event logevent.LogEvent) (err error) {
	var (
		conn redis.Conn
		raw  []byte
		key  string
	)

	if raw, err = event.MarshalJSON(); err != nil {
		log.Errorf("event Marshal failed: %v", event)
		return
	}
	key = event.Format(self.Key)

	conn = self.conns.Get()
	defer conn.Close()

	switch self.DataType {
	case "list":
		fmt.Println(string(key) + string(raw))
		_, err = conn.Do("rpush", key, raw)
		fmt.Println(err)
	case "channel":
		_, err = conn.Do("publish", key, raw)
	default:
		err = errors.New("unknown DataType: " + self.DataType)
	}

	return
}

func initRedisPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     200 * 2,
		MaxActive:   200 * 2,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password == "" {
				return c, err
			}
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
