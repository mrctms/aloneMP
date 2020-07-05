package clients

import "aloneMP/senders"

type Clienter interface {
	SetSender(sender senders.Sender)
	Run(source string)
}
