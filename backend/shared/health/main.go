package health

import (
	repositories "git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/database"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/constants"
	"git.asseco-see.hr/asseco-hr-voice/evil/go-chassis/v2/pkg/shared/helpers"
	"github.com/Azure/go-amqp"
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
	vab := ValidateAmqpBus()
	if vab != nil {
		amqpRes = vab.Error()
	}
	return &Response{
		Database: database,
		Amqp:     amqpRes,
	}
}

func ValidateDatabase() error {
	conn, err := repositories.NewConnection(
		helpers.GetEnvString(constants.DatabaseDriver, ""),
		helpers.GetEnvString(constants.DatabaseUser, ""),
		helpers.GetEnvString(constants.DatabasePass, ""),
		helpers.GetEnvString(constants.DatabaseHost, ""),
		helpers.GetEnvString(constants.DatabasePort, ""),
		helpers.GetEnvString(constants.DatabaseName, ""),
		helpers.GetEnvBool(constants.DebugDatabase, false),
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

func ValidateAmqpBus() error {
	client, err := amqp.Dial("amqp://"+helpers.GetEnvString(constants.BrokerHost, "")+":"+helpers.GetEnvString(constants.BrokerPort, ""),
		amqp.ConnSASLPlain(helpers.GetEnvString(constants.BrokerUser, ""), helpers.GetEnvString(constants.BrokerPass, "")),
	)
	if err != nil {
		return err
	}
	return client.Close()
}
