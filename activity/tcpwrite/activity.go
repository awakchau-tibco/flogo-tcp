package tcpwrite

import (
	"fmt"
	"net"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

// Activity ...
type Activity struct {
	settings   *Settings
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
	activity := &Activity{}
	ctx.Logger().Debug("Dialing connection...")
	if s.Network == "" {
		s.Network = "tcp"
	}
	activity.connection, err = net.Dial(s.Network, fmt.Sprintf("%s:%s", s.Host, s.Port))
	if err != nil {
		ctx.Logger().Errorf("Unable to dial the connection! %s", err.Error())
		return nil, err
	}
	ctx.Logger().Debug("Connection is now open")
	activity.settings = s
	return activity, nil
}

// Cleanup ...
func (a *Activity) Cleanup() error {
	fmt.Println("Cleaning up...")
	err := a.connection.Close()
	if err != nil {
		return err
	}
	return nil
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
	output := &Output{}
	if a.settings.WriteTimeoutMs != 0 {
		deadline := time.Now().Add(time.Millisecond * time.Duration(a.settings.WriteTimeoutMs))
		a.connection.SetWriteDeadline(deadline)
	}
	output.BytesWritten, err = a.connection.Write([]byte(message))
	if err != nil {
		ctx.Logger().Errorf("Unable to write the data! %s", err.Error())
		return false, err
	}

	ctx.SetOutputObject(output)
	ctx.Logger().Infof("Written %d bytes", output.BytesWritten)
	return true, nil
}
