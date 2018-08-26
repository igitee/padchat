package padchat_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tuotoo/padchat"
)

func TestMkAtContent(t *testing.T) {
	req := &padchat.SendMsgReq{
		Content: "test",
	}
	padchat.MkAtContent(req)
	assert.Equal(t, "test", req.Content)
	req.AtList = []string{"1", "2", "3"}
	padchat.MkAtContent(req)
	assert.Equal(t, "@@@\ntest", req.Content)
	req.Content = "@@\ntest"
	padchat.MkAtContent(req)
	assert.Equal(t, "@\n@@\ntest", req.Content)
}
