package tcpwrite

import (
	"testing"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/support/test"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {

	ref := activity.GetRef(&Activity{})
	act := activity.Get(ref)

	assert.NotNil(t, act)
}

func TestEval(t *testing.T) {
	settings := Settings{
		Port:      "8888",
		Delimiter: "\n",
	}
	initContext := test.NewActivityInitContext(settings, nil)
	act, err := New(initContext)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())

	aInput := &Input{
		StringData: []byte("Message 1\nMessage 2;Message 3"),
	}

	tc.SetInputObject(aInput)

	done, _ := act.Eval(tc)

	assert.True(t, done)
	aOutput := &Output{}
	err = tc.GetOutputObject(aOutput)
	assert.Nil(t, err)
	assert.Greater(t, aOutput.BytesWritten, 0)
}
