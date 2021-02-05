package tcpwrite

import "github.com/project-flogo/core/data/coerce"

// Settings ...
type Settings struct {
	Network        string `md:"network"`        // The network type
	Host           string `md:"host"`           // The host name or IP for TCP server.
	Port           string `md:"port,required"`  // The port to listen on
	WriteTimeoutMs int64  `md:"writeTimeoutMs"` // Write timeout for tcp write operation in milliseconds
}

// ToMap ...
func (i *Settings) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"network":        i.Network,
		"host":           i.Host,
		"port":           i.Port,
		"writeTimeoutMs": i.WriteTimeoutMs,
	}
}

// FromMap ...
func (i *Settings) FromMap(values map[string]interface{}) error {
	var err error
	i.Network, err = coerce.ToString(values["network"])
	if err != nil {
		return err
	}
	i.Host, err = coerce.ToString(values["host"])
	if err != nil {
		return err
	}
	i.Port, err = coerce.ToString(values["port"])
	if err != nil {
		return err
	}
	i.WriteTimeoutMs, err = coerce.ToInt64(values["writeTimeoutMs"])
	if err != nil {
		return err
	}
	return nil
}

// Input ...
type Input struct {
	StringData string `md:"stringData,required"`
	Delimiter  string `md:"delimiter"` // Data delimiter for read and write
}

// ToMap ...
func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"stringData": i.StringData,
		"delimiter":  i.Delimiter,
	}
}

// FromMap ...
func (i *Input) FromMap(values map[string]interface{}) error {
	var err error
	i.StringData, err = coerce.ToString(values["stringData"])
	if err != nil {
		return err
	}
	i.Delimiter, err = coerce.ToString(values["delimiter"])
	if err != nil {
		return err
	}
	return nil
}

// Output ...
type Output struct {
	BytesWritten int `md:"bytesWritten"`
}

// ToMap ...
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"bytesWritten": o.BytesWritten,
	}
}

// FromMap ...
func (o *Output) FromMap(values map[string]interface{}) error {
	var err error
	o.BytesWritten, err = coerce.ToInt(values["bytesWritten"])
	if err != nil {
		return err
	}
	return nil
}
