package paymentRepositories

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/korvised/ilog-shop/modules/inventory"
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
		log.Printf("Payment Error: DockedPlayerMoney: %s\n", err.Error())
		return errors.New("error: docked player money failed")
	}

	log.Println("Pushed message to player")

	return nil
}

func (r *paymentRepository) RollbackTransaction(_ context.Context, req *player.RollbackPlayerTransactionReq) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Payment Error: RollbackTransaction: %s\n", err.Error())
		return errors.New("error: rollback player transaction failed")
	}

	log.Printf("Payment Info: RollbackTransaction: %s\n", string(reqInBytes))

	if err = queue.PushMessageWithKeyToQueue(
		[]string{r.cfg.Kafka.Url},
		r.cfg.Kafka.ApiKey,
		r.cfg.Kafka.Secret,
		"player",
		"rollback_buy",
		reqInBytes,
	); err != nil {
		log.Printf("Error: RollbackTransaction: %s\n", err.Error())
		return errors.New("error: rollback player transaction failed")
	}

	return nil
}

func (r *paymentRepository) AddPlayItem(_ context.Context, req *inventory.UpdateInventoryReq) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Payment Error: AddPlayItem: %s\n", err.Error())
		return errors.New("error: add player item failed")
	}

	log.Printf("Payment Info: AddPlayItem: %s\n", string(reqInBytes))

	if err = queue.PushMessageWithKeyToQueue(
		[]string{r.cfg.Kafka.Url},
		r.cfg.Kafka.ApiKey,
		r.cfg.Kafka.Secret,
		"inventory",
		"buy",
		reqInBytes,
	); err != nil {
		log.Printf("Error: AddPlayItem: %s\n", err.Error())
		return errors.New("error: add player item failed")
	}

	return nil
}

func (r *paymentRepository) RollbackAddPlayItem(_ context.Context, req *inventory.RollbackPlayerInventoryReq) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Payment Error: RollbackAddPlayItem: %s\n", err.Error())
		return errors.New("error: rollback add player item failed")
	}

	log.Printf("Payment Info: RollbackAddPlayItem: %s\n", string(reqInBytes))

	if err = queue.PushMessageWithKeyToQueue(
		[]string{r.cfg.Kafka.Url},
		r.cfg.Kafka.ApiKey,
		r.cfg.Kafka.Secret,
		"inventory",
		"rollback_buy",
		reqInBytes,
	); err != nil {
		log.Printf("Error: RollbackAddPlayItem: %s\n", err.Error())
		return errors.New("error: rollback add player item failed")
	}

	return nil
}

func (r *paymentRepository) RemovePlayItem(_ context.Context, req *inventory.UpdateInventoryReq) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Payment Error: RemovePlayItem: %s\n", err.Error())
		return errors.New("error: remove player item failed")
	}

	log.Printf("Payment Info: RemovePlayItem: %s\n", string(reqInBytes))

	if err = queue.PushMessageWithKeyToQueue(
		[]string{r.cfg.Kafka.Url},
		r.cfg.Kafka.ApiKey,
		r.cfg.Kafka.Secret,
		"inventory",
		"sell",
		reqInBytes,
	); err != nil {
		log.Printf("Error: RemovePlayItem: %s\n", err.Error())
		return errors.New("error: remove player item failed")
	}

	return nil
}

func (r *paymentRepository) RollbackRemovePlayItem(_ context.Context, req *inventory.RollbackPlayerInventoryReq) error {
	reqInBytes, err := json.Marshal(req)
	if err != nil {
		log.Printf("Payment Error: RollbackRemovePlayItem: %s\n", err.Error())
		return errors.New("error: rollback remove player item failed")
	}

	log.Printf("Payment Info: RollbackRemovePlayItem: %s\n", string(reqInBytes))

	if err = queue.PushMessageWithKeyToQueue(
		[]string{r.cfg.Kafka.Url},
		r.cfg.Kafka.ApiKey,
		r.cfg.Kafka.Secret,
		"inventory",
		"rollback_sell",
		reqInBytes,
	); err != nil {
		log.Printf("Error: RollbackRemovePlayItem: %s\n", err.Error())
		return errors.New("error: rollback remove player item failed")
	}

	return nil
}
