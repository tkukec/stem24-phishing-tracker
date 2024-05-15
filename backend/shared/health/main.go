package health

import (
	"github.com/andrezz-b/stem24-phishing-tracker/shared/constants"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/database"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/runtimebag"
	"net/http"
)

const OkString = "OK"

type Response struct {
	Database string `json:"database"`
	Amqp     string `json:"amqp"`
}

func (r *Response) StatusCode() int {
	if r.Database != OkString || r.Amqp != OkString {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func Validate() *Response {
	database := OkString
	amqpRes := OkString
	err := ValidateDatabase()
	if err != nil {
		database = err.Error()
	}

	return &Response{
		Database: database,
		Amqp:     amqpRes,
	}
}

func ValidateDatabase() error {
	conn, err := database.NewConnection(
		runtimebag.GetEnvString(constants.DatabaseDriver, ""),
		runtimebag.GetEnvString(constants.DatabaseUser, ""),
		runtimebag.GetEnvString(constants.DatabasePass, ""),
		runtimebag.GetEnvString(constants.DatabaseHost, ""),
		runtimebag.GetEnvString(constants.DatabasePort, ""),
		runtimebag.GetEnvString(constants.DatabaseName, ""),
		runtimebag.GetEnvBool(constants.DebugDatabase, false),
		nil,
	)
	if err != nil {
		return nil
	}
	db, err := conn.GetConnection().DB()
	if err != nil {
		return nil
	}
	return db.Close()
}
