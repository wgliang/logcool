// Output-plug: outputemail
// The plug's function take the event-data into the standard-output.
package outputemail

import (
	"strconv"
	"strings"

	"github.com/wgliang/logcool/utils"
	"gopkg.in/gomail.v2"
)

const (
	ModuleName = "email"
)

// Define outputemail' config.
type OutputConfig struct {
	utils.OutputConfig

	Server   string
	From     string
	Password string
	To       []string
	Cc       string
}

func init() {
	utils.RegistOutputHandler(ModuleName, InitHandler)
}

// Init outputemail Handler.
func InitHandler(confraw *utils.ConfigRaw) (retconf utils.TypeOutputConfig, err error) {
	conf := OutputConfig{
		OutputConfig: utils.OutputConfig{
			CommonConfig: utils.CommonConfig{
				Type: ModuleName,
			},
		},
	}
	if err = utils.ReflectConfig(confraw, &conf); err != nil {
		return
	}

	retconf = &conf
	return
}

// Input's event,and this is the main function of output.
func (oc *OutputConfig) Event(event utils.LogEvent) (err error) {

	return oc.Send(event)
}

func (oc *OutputConfig) Send(event utils.LogEvent) error {
	raw, err := event.MarshalIndent()
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", oc.From)
	m.SetHeader("To", oc.To...)
	m.SetAddressHeader("Cc", oc.Cc, "all")
	m.SetHeader("Subject", event.Message)
	m.SetBody("text/html", string(raw))
	// m.Attach("/tmp/log/detail.log")

	server := strings.Split(oc.Server, ":")
	username := strings.Split(oc.From, "@")
	port, err := strconv.Atoi(server[1])
	if err != nil {
		return err
	}

	d := gomail.NewDialer(server[0], port, username[0], oc.Password)

	// Send the email.
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
