package enc

import (
	"fmt"
	"github.com/myl7/datasize"
	"os"
	"sync"
)

type Config struct {
	InputMaxsize uint64
}

type ErrInvalidConfig struct {
	Field  string
	ErrVal string
}

func (e ErrInvalidConfig) Error() string {
	return fmt.Sprintf("invalid config: %s = %s", e.Field, e.ErrVal)
}

func NewConfig() (Config, error) {
	c := Config{}

	s := os.Getenv("BROTLI_INPUT_MAXSIZE")
	var b datasize.ByteSize
	err := b.UnmarshalText([]byte(s))
	if err != nil {
		return Config{}, ErrInvalidConfig{Field: "BROTLI_INPUT_MAXSIZE", ErrVal: s}
	}

	c.InputMaxsize = b.Bytes()

	return c, err
}

var config *Config
var configLock sync.RWMutex

func GetConfig() Config {
	isRead := true
	configLock.RLock()
	defer func() {
		if isRead {
			configLock.RUnlock()
		} else {
			configLock.Unlock()
		}
	}()

	if config == nil {
		isRead = false
		configLock.Lock()
		c, err := NewConfig()
		if err != nil {
			panic(err)
		}

		config = &c
		return c
	} else {
		return *config
	}
}
