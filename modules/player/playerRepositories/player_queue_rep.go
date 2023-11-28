package playerRepositories

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/korvised/ilog-shop/modules/payment"
	"github.com/korvised/ilog-shop/pkg/queue"
	"log"
)

func (r *playerRepository) DockedPlayerMoneyRes(_ context.Context, req *payment.PaymentTransferRes) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: DockedPlayerMoneyRes: %s\n", err.Error())
		return errors.New("error: docked player money failed")
	}

	if err = queue.PushMessageWithKeyToQueue(
		[]string{r.cfg.Kafka.Url},
		r.cfg.Kafka.ApiKey,
		r.cfg.Kafka.Secret,
		"payment",
		"buy",
		reqInBytes,
	); err != nil {
		log.Printf("Error: DockedPlayerMoneyRes: %s\n", err.Error())
		return errors.New("error: docked player money failed")
	}

	return nil
}
