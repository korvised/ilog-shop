package paymentRepositories

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/korvised/ilog-shop/modules/player"
	"github.com/korvised/ilog-shop/pkg/queue"
	"log"
)

func (r *paymentRepository) DockedPlayerMoney(_ context.Context, req *player.CreatePlayerTransactionReq) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: DockedPlayerMoney: %s\n", err.Error())
		return errors.New("error: docked player money failed")
	}

	if err = queue.PushMessageWithKeyToQueue(
		[]string{r.cfg.Kafka.Url},
		r.cfg.Kafka.ApiKey,
		r.cfg.Kafka.Secret,
		"player",
		"buy",
		reqInBytes,
	); err != nil {
		log.Printf("Error: DockedPlayerMoney: %s\n", err.Error())
		return errors.New("error: docked player money failed")
	}

	log.Println("Pushed message to player")

	return nil
}

func (r *paymentRepository) RollbackTransaction(_ context.Context, req *player.RollbackPlayerTransactionReq) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error: RollbackTransaction: %s\n", err.Error())
		return errors.New("error: rollback player money failed")
	}

	if err = queue.PushMessageWithKeyToQueue(
		[]string{r.cfg.Kafka.Url},
		r.cfg.Kafka.ApiKey,
		r.cfg.Kafka.Secret,
		"player",
		"rtransaction",
		reqInBytes,
	); err != nil {
		log.Printf("Error: RollbackTransaction: %s\n", err.Error())
		return errors.New("error: rollback player money failed")
	}

	return nil
}
