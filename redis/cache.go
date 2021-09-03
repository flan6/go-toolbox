package redis

import (
	"fmt"

	redigo "github.com/garyburd/redigo/redis"
)

// Args is a helper for constructing command arguments from structured values.
type Args []interface{}

// Add returns the result of appending value to args.
func (args Args) Add(value ...interface{}) Args {
	return append(args, value...)
}

// Client is the structure used to create a redis connection client
type Client struct {
	Host   string
	Port   int
	DB     int
	Prefix string
	conn   *redigo.Conn
}

// DefaultClient is a default client to connect to local redis
var DefaultClient = Client{
	Host:   "localhost",
	DB:     0,
	Port:   6379,
	Prefix: "",
}

// Conn returns a redis connection to execute commands
func (c *Client) Conn() (conn redigo.Conn) {
	if c.conn != nil {
		return *c.conn
	}

	conn, err := redigo.Dial("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port), redigo.DialDatabase(c.DB))
	if err != nil {
		panic(err.Error())
	}

	c.conn = &conn
	return *c.conn
}

// Set the string value of a key
func (c Client) Set(key, value string, expirationSeconds ...int) (err error) {
	key = c.Prefix + key

	_, err = c.Conn().Do("SET", key, value)
	if err != nil {
		return
	}

	if len(expirationSeconds) > 0 {
		_, err = c.Conn().Do("EXPIRE", key, expirationSeconds[0])
	}

	return
}

// Get the value of a key
func (c Client) Get(key string) (value string, err error) {
	return redigo.String(c.Conn().Do("GET", c.Prefix+key))
}

// MustGet the value of a key and you can check for a boolean returned
func (c Client) MustGet(key string) (value string, ok bool) {
	var err error
	value, err = c.Get(key)
	if err != nil || value == "" {
		return "", false
	}

	return value, true
}

// Delete a key
func (c Client) Delete(key string) (err error) {
	_, err = c.Conn().Do("DEL", c.Prefix+key)
	return
}

// Delete todas as chaves onde contém o pattern localizado
func (c Client) DeleteLike(pattern string) (err error) {
	iter := 0
	for {
		arr, err := redigo.Values(c.Conn().Do("SCAN", iter, "MATCH", "*"+pattern+"*"))
		if err != nil {
			return fmt.Errorf("error retrieving '%s' keys", c.Prefix+pattern)
		}

		iter, _ = redigo.Int(arr[0], nil)
		keys, _ := redigo.Strings(arr[1], nil)

		for _, key := range keys {
			_, err = c.Conn().Do("DEL", key)

			if err != nil {
				return err
			}
		}

		if iter == 0 {
			break
		}
	}

	return nil
}

//Do Abre uma conexão com o Redis, executa o comando e depois a fecha
func (c Client) Send(comando string, args ...interface{}) (interface{}, error) {
	value := c.Conn().Send(comando, args...)
	if value != nil {
		return nil, value
	}
	return value, nil
}

//**** HM

/*// HMGet the value of a key
func (c Client) HMGet(key string) (values []string, err error) {
	return redigo.Strings(c.Conn().Do("HMGET", c.Prefix+key))
}

// MustGet the value of a key and you can check for a boolean returned
func (c Client) HMMustGet(key string) (values []string, ok bool) {
	var err error
	values, err = c.HMGet(key)
	if err != nil || len(values) == 0 {
		return []string{}, false
	}

	return values, true
}


// HMSet the string value of a key
func (c Client) HMSet(key string, expirationSeconds int,
	values ...interface{}) (err error) {
	key = c.Prefix + key

	valuesArray := Args{}.Add(key).Add(values...)

	_, err = c.Conn().Do("HMSET", valuesArray)
	if err != nil {
		return
	}

	if expirationSeconds > 0 {
		_, err = c.Conn().Do("EXPIRE", key, expirationSeconds)
	}

	return
}

// HMDelete a key
func (c Client) HMDelete(key string) (err error) {
	_, err = c.Conn().Do("HMDEL", c.Prefix+key)
	return
}*/
