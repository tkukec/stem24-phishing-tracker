package runtimebag

import (
	"log"
	"os"
	"strconv"
)

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// NewBagWithPreloadedEnvs same as NewBag but reads env variables and saves them into the bag
func NewBagWithPreloadedEnvs() *Bag {
	bag := NewBag()
	for _, pair := range os.Environ() {
		parts := strings.Split(pair, "=")
		bag.SetEnvValue(parts[0], strings.Join(parts[1:], "="))
	}
	return bag
}

// NewBag a quick storage to share for simple values. NOTE: THIS IS NOT A REPLACEMENT FOR REDIS AND SHOULD NOT BE USED AS SUCH
func NewBag() *Bag {
	return &Bag{
		values:    make(map[string]interface{}),
		envValues: make(map[string]interface{}),
	}
}

type Bag struct {
	values    map[string]interface{}
	envValues map[string]interface{}
	rm        sync.RWMutex
}

func (b *Bag) AddValue(key string, value interface{}) {
	b.rm.Lock()
	defer b.rm.Unlock()
	b.values[key] = value
}

func (b *Bag) GetAllValues() map[string]interface{} {
	return b.values
}

func (b *Bag) GetAllEnvValues() map[string]interface{} {
	return b.envValues
}

func (b *Bag) SetEnvValue(key string, value interface{}) error {

	if _, ok := b.envValues[key]; ok {
		return fmt.Errorf("value already set for key %s. You can not change an ENV value", key)
	}

	b.envValues[key] = value
	return nil
}

func (b *Bag) InvalidateValue(key string) {
	delete(b.values, key)
}

func (b *Bag) GetValue(key string) interface{} {
	return b.values[key]
}

func (b *Bag) GetString(key string) (string, bool, error) {
	b.rm.Lock()
	defer b.rm.Unlock()
	if value, ok := b.values[key]; ok {
		switch cast := value.(type) {
		case string:
			return cast, true, nil
		default:
			return "", false, fmt.Errorf("value is not string. correct value %s", reflect.TypeOf(value))
		}
	}
	return "", false, fmt.Errorf("value for %s not found", key)
}

func (b *Bag) GetInteger(key string) (int, bool, error) {
	if value, ok := b.values[key]; ok {
		switch cast := value.(type) {
		case int:
			return cast, true, nil
		default:
			return 0, false, fmt.Errorf("value is not integer. correct value %s", reflect.TypeOf(value))
		}
	}
	return 0, false, fmt.Errorf("value for %s not found", key)
}

func (b *Bag) GetStringSlice(key string) ([]string, bool, error) {
	if value, ok := b.values[key]; ok {
		switch cast := value.(type) {
		case []string:
			return cast, true, nil
		default:
			return nil, false, fmt.Errorf("value is not string slice. correct value %s", reflect.TypeOf(value))
		}
	}
	return nil, false, fmt.Errorf("value for %s not found", key)
}

func (b *Bag) GetIntegerSlice(key string) ([]int, bool, error) {
	if value, ok := b.values[key]; ok {
		switch cast := value.(type) {
		case []int:
			return cast, true, nil
		default:
			return nil, false, fmt.Errorf("value is not integer slice. correct value %s", reflect.TypeOf(value))
		}
	}
	return nil, false, fmt.Errorf("value for %s not found", key)
}

func (b *Bag) GetEnvValue(key string) interface{} {
	return b.envValues[key]
}

func (b *Bag) GetEnvString(key string) (string, bool, error) {
	if value, ok := b.envValues[key]; ok {
		switch cast := value.(type) {
		case string:
			return cast, true, nil
		default:
			return "", false, fmt.Errorf("value is not string. correct value %s", reflect.TypeOf(value))
		}
	}
	return "", false, fmt.Errorf("value for %s not found", key)
}

func (b *Bag) GetEnvInteger(key string) (int, bool, error) {
	if value, ok := b.envValues[key]; ok {
		switch cast := value.(type) {
		case int:
			return cast, true, nil
		default:
			return 0, false, fmt.Errorf("value is not integer. correct value %s", reflect.TypeOf(value))
		}
	}
	return 0, false, fmt.Errorf("value for %s not found", key)
}

func (b *Bag) GetEnvStringSlice(key string) ([]string, bool, error) {
	if value, ok := b.envValues[key]; ok {
		switch cast := value.(type) {
		case []string:
			return cast, true, nil
		default:
			return nil, false, fmt.Errorf("value is not string slice. correct value %s", reflect.TypeOf(value))
		}
	}
	return nil, false, fmt.Errorf("value for %s not found", key)
}

func (b *Bag) GetEnvIntegerSlice(key string) ([]int, bool, error) {
	if value, ok := b.envValues[key]; ok {
		switch cast := value.(type) {
		case []int:
			return cast, true, nil
		default:
			return nil, false, fmt.Errorf("value is not integer slice. correct value %s", reflect.TypeOf(value))
		}
	}
	return nil, false, fmt.Errorf("value for %s not found", key)
}

func GetEnvString(key, defaultValue string) string {
	if x := os.Getenv(key); x != "" {
		return x
	}
	return defaultValue
}

func GetEnvInt(key string, defaultValue int64) int64 {
	if x := os.Getenv(key); x != "" {
		v, err := strconv.ParseInt(x, 10, 64)
		if err != nil {
			log.Panicf("Can not convert string to int64 %s", err.Error())
		}
		return int64(v)
	}
	return defaultValue
}

func GetEnvBool(key string, defaultValue bool) bool {
	if x := os.Getenv(key); x != "" {
		v, err := strconv.ParseBool(x)
		if err != nil {
			log.Panicf("Can not convert string to bool %s", err.Error())
		}
		return v
	}
	return defaultValue
}
