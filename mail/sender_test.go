package mail

import (
	"testing"

	"github.com/corey888773/golang-course/util"
	"github.com/stretchr/testify/require"
)

func TestSendingEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "a test email"
	content := `
	<h1>Hello World!</h1>
	<p>This is a test message</p>
	`
	to := []string{"piotrdropii4@gmail.com"}
	attachFiles := []string{"../Makefile"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
