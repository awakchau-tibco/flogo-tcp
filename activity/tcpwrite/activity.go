package tcpwrite

import (
	"fmt"
	"net"
	"strings"
	"unicode/utf8"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

// Activity ...
type Activity struct {
	connection net.Conn
}

func init() {
	_ = activity.Register(&Activity{}, New)
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// New ...
func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}
	act := &Activity{}
	ctx.Logger().Debug("Dialing connection...")
	if s.Network == "" {
		s.Network = "tcp"
	}
	act.connection, err = net.Dial(s.Network, fmt.Sprintf("%s:%s", s.Host, s.Port))
	if err != nil {
		ctx.Logger().Errorf("Unable to dial the connection! %s", err.Error())
		return nil, err
	}
	ctx.Logger().Debug("Connection is now open")
	return act, nil
}

// Eval ...
func (a *Activity) Eval(context activity.Context) (done bool, err error) {
	context.Logger().Debug("Executing TCP Write activity")
	input := &Input{}
	err = context.GetInputObject(input)
	if err != nil {
		return false, err
	}
	output := &Output{}
	message := input.StringData
	if input.Delimiter != "" {
		r, _ := utf8.DecodeRuneInString(input.Delimiter)
		delimiter := byte(r)
		message = input.StringData[:strings.IndexByte(message, delimiter)]
	}
	defer a.connection.Close()
	output.BytesWritten, err = a.connection.Write([]byte(message))
	if err != nil {
		context.Logger().Errorf("Unable to write the data! %s", err.Error())
		return false, err
	}
	context.SetOutputObject(output)
	context.Logger().Debug("TCP Write activity completed")
	return true, nil
}
