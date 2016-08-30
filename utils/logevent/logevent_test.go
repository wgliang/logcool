package logevent

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_AddTag(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

	le := &LogEvent{
		Timestamp: time.Now(),
		Message:   "message",
		Tags:      []string{"frg", "grbhrt"},
		Extra:     make(map[string]interface{}),
	}
	le.AddTag("vftb")
	fmt.Println(le.Tags)
}

func Test_MarshalJSON(t *testing.T) {

	le := &LogEvent{
		Timestamp: time.Now(),
		Message:   "message",
		Tags:      []string{"frg", "grbhrt"},
		Extra:     make(map[string]interface{}),
	}
	data, err := le.MarshalJSON()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(data))
}

func Test_MarshalIndent(t *testing.T) {
	le := &LogEvent{
		Timestamp: time.Now(),
		Message:   "message",
		Tags:      []string{"frg", "grbhrt"},
		Extra:     make(map[string]interface{}),
	}
	data, err := le.MarshalIndent()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))
}

func Test_Get(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

	le := &LogEvent{
		Timestamp: time.Now(),
		Message:   "message",
		Tags:      []string{"frg", "grbhrt"},
		Extra:     make(map[string]interface{}),
	}
	data := le.Get("@timestamp")
	strs := string(data.(time.Time).UTC().Format(timeFormat))
	fmt.Println(strs)

	data = le.Get("message")
	assert.Equal("message", data)

}

func Test_GetString(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

	le := &LogEvent{
		Timestamp: time.Now(),
		Message:   "message",
		Tags:      []string{"frg", "grbhrt"},
		Extra:     make(map[string]interface{}),
	}
	data := le.GetString("@timestamp")
	fmt.Println(data)

	data = le.GetString("message")
	assert.Equal("message", data)

	data = le.GetString("extra")
	fmt.Println(data)
}

func Test_FormatWithEnv(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

	key := "TESTENV"

	originenv := os.Getenv(key)
	defer func() {
		os.Setenv(key, originenv)
	}()

	err := os.Setenv(key, "Testing ENV")
	assert.NoError(err)

	out := FormatWithEnv("prefix %{TESTENV} suffix")
	assert.Equal("prefix Testing ENV suffix", out)
}

func Test_FormatWithTime(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

	out := FormatWithTime("prefix %{+2006-01-02} suffix")
	nowdatestring := time.Now().Format("2006-01-02")
	assert.Equal("prefix "+nowdatestring+" suffix", out)
}

func Test_Format(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)

	logevent := LogEvent{
		Timestamp: time.Now(),
		Message:   "Test Message",
		Extra: map[string]interface{}{
			"int":    123,
			"float":  1.23,
			"string": "Test String",
			"time":   time.Now(),
		},
	}

	out := logevent.Format("%{message}")
	assert.Equal("Test Message", out)

	out = logevent.Format("%{@timestamp}")
	assert.NotEmpty(out)
	assert.NotEqual("%{@timestamp}", out)

	out = logevent.Format("%{int}")
	assert.Equal("123", out)

	out = logevent.Format("%{float}")
	assert.Equal("1.23", out)

	out = logevent.Format("%{string}")
	assert.Equal("Test String", out)

	out = logevent.Format("time string %{+2006-01-02}")
	nowdatestring := time.Now().Format("2006-01-02")
	assert.Equal("time string "+nowdatestring, out)

	out = logevent.Format("%{null}")
	assert.Equal("%{null}", out)

	logevent.AddTag("tag1", "tag2", "tag3")
	assert.Len(logevent.Tags, 3)
	assert.Contains(logevent.Tags, "tag1")

	logevent.AddTag("tag1", "tag%{int}")
	assert.Len(logevent.Tags, 4)
	assert.Contains(logevent.Tags, "tag123")
}
