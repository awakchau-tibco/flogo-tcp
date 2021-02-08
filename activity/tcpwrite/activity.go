package tcpwrite

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
)

var logger log.Logger

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
	settings := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), settings, true)
	if err != nil {
		return nil, err
	}
	activity := &Activity{}
	logger = ctx.Logger()
	if settings.Network == "" {
		settings.Network = "tcp"
	}
	logger.Debug(fmt.Sprintf("Dialing connection using %s network...", settings.Network))
	activity.connection, err = net.Dial(settings.Network, fmt.Sprintf("%s:%s", settings.Host, settings.Port))
	if err != nil {
		logger.Errorf("Unable to dial the connection! Caused by %s", err.Error())
		return nil, err
	}
	logger.Debug("Connection is now open")
	activity.settings = settings
	return activity, nil
}

// Cleanup ...
func (a *Activity) Cleanup() error {
	logger.Info("Closing connection")
	err := a.connection.Close()
	if err != nil {
		return err
	}
	return nil
}

// Eval ...
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
	logger.Debug("Executing TCP Write activity")
	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return false, err
	}
	message := input.StringData
	if len(a.settings.Delimiter) > 0 {
		message = input.StringData[:bytes.Index(input.StringData, []byte(a.settings.Delimiter))]
	}
	output := &Output{}
	if a.settings.WriteTimeoutMs != 0 {
		deadline := time.Now().Add(time.Millisecond * time.Duration(a.settings.WriteTimeoutMs))
		a.connection.SetWriteDeadline(deadline)
	}
	output.BytesWritten, err = a.connection.Write(message)
	if err != nil {
		logger.Errorf("Unable to write the data! %s", err.Error())
		return false, err
	}
	ctx.SetOutputObject(output)
	logger.Infof("Written %d bytes", output.BytesWritten)
	return true, nil
}
