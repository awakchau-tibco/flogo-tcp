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
	settings *Settings
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
	act.settings = s
	return act, nil
}

// Eval ...
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
	ctx.Logger().Debug("Executing TCP Write activity")
	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return false, err
	}
	message := input.StringData
	if input.Delimiter != "" {
		r, _ := utf8.DecodeRuneInString(input.Delimiter)
		delimiter := byte(r)
		message = input.StringData[:strings.IndexByte(message, delimiter)]
	}

	conn, err := net.Dial(a.settings.Network, fmt.Sprintf("%s:%s", a.settings.Host, a.settings.Port))
	if err != nil {
		ctx.Logger().Errorf("Unable to dial the connection! %s", err.Error())
		return false, err
	}
	defer conn.Close()
	ctx.Logger().Debug("Connection is now open")
	output := &Output{}

	output.BytesWritten, err = conn.Write([]byte(message + input.Delimiter))
	if err != nil {
		ctx.Logger().Errorf("Unable to write the data! %s", err.Error())
		return false, err
	}

	ctx.SetOutputObject(output)
	ctx.Logger().Infof("Written %d bytes", output.BytesWritten)
	return true, nil
}
