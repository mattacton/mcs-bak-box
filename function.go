// Package p contains a Pub/Sub Cloud Function.
package p

import (
	"context"
	"log"

	"github.com/mattacton/mcs-bak-box/mcsbox"
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

func MCSBakBoxPubSub(ctx context.Context, m PubSubMessage) error {
	log.Println(string(m.Data))
	mcsbox.BackupMCSBox()
	return nil
}
