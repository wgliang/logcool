// Output-plug: outputredis
// The plug's function take the event-data into redis.
package outputredis

import (
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"

	"github.com/wgliang/logcool/utils"
)

const (
	ModuleName = "redis"
)

// Define outputredis' config.
type OutputConfig struct {
	utils.OutputConfig
	Key               string `json:"key"`
	Host              string `json:"host"`
	Password          string `json:"password"`
	DataType          string `json:"data_type"`
	Timeout           int    `json:"timeout"`
	ReconnectInterval int    `json:"reconnect_interval"`

	conns  *redis.Pool
	evchan chan utils.LogEvent
}

// Init outputredis Handler.
func InitHandler(confraw *utils.ConfigRaw) (retconf utils.TypeOutputConfig, err error) {
	conf := OutputConfig{
		OutputConfig: utils.OutputConfig{
			CommonConfig: utils.CommonConfig{
				Type: ModuleName,
			},
		},

		evchan: make(chan utils.LogEvent),
	}
	if err = utils.ReflectConfig(confraw, &conf); err != nil {
		return
	}
	conf.conns = initRedisPool(conf.Host, conf.Password)

	go conf.loopEvent()

	retconf = &conf
	return
}

// Input's event,and this is the main function of output.
func (oc *OutputConfig) Event(event utils.LogEvent) (err error) {
	oc.evchan <- event
	return
}

func (oc *OutputConfig) loopEvent() (err error) {
	for {
		event := <-oc.evchan
		oc.sendEvent(event)
	}
}

func (oc *OutputConfig) sendEvent(event utils.LogEvent) (err error) {
	var (
		conn redis.Conn
		raw  []byte
		key  string
	)

	if raw, err = event.MarshalJSON(); err != nil {
		log.Errorf("event Marshal failed: %v", event)
		return
	}
	key = event.Format(oc.Key)

	// get a connection from pool
	conn = oc.conns.Get()
	defer conn.Close()

	switch oc.DataType {
	case "list":
		_, err = conn.Do("rpush", key, raw)
	case "set":
		_, err = conn.Do("set", key, raw)
	case "channel":
		_, err = conn.Do("publish", key, raw)
	case "append":
		_, err = conn.Do("publish", key, raw)
	default:
		err = errors.New("unknown DataType: " + oc.DataType)
	}

	return
}

// init a redis connection pool.
func initRedisPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     200 * 2,
		MaxActive:   200 * 2,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			// test tcp.
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
