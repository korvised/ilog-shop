package playerHandlers

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/player"
	"github.com/korvised/ilog-shop/modules/player/playerUsecases"
	"github.com/korvised/ilog-shop/pkg/queue"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type (
	PlayerQueueHandlerService interface {
		playerConsumer(c context.Context) (sarama.PartitionConsumer, error)
		DockedPlayerMoney()
		RollbackPlayerTransaction()
	}

	playerQueueHandler struct {
		cfg           *config.Config
		playerUsecase playerUsecases.PlayerUsecaseService
	}
)

func NewPlayerQueueHandler(cfg *config.Config, playerUsecase playerUsecases.PlayerUsecaseService) PlayerQueueHandlerService {
	return &playerQueueHandler{cfg, playerUsecase}
}

func (h *playerQueueHandler) playerConsumer(c context.Context) (sarama.PartitionConsumer, error) {
	worker, err := queue.ConnectConsumer([]string{h.cfg.Kafka.Url}, h.cfg.Kafka.ApiKey, h.cfg.Kafka.Secret)
	if err != nil {
		return nil, err
	}

	offset, err := h.playerUsecase.GetOffset(c)
	if err != nil {
		return nil, err
	}
	consumer, err := worker.ConsumePartition("player", 0, offset)
	if err != nil {
		log.Printf("Error: PlayerConsumer 1: %v \n", err)
		consumer, err = worker.ConsumePartition("player", 0, offset)
		if err != nil {
			log.Printf("Error: PlayerConsumer 2 : %v \n", err)
			return nil, err
		}
	}

	return consumer, nil
}

func (h *playerQueueHandler) DockedPlayerMoney() {
	ctx := context.Background()

	consumer, err := h.playerConsumer(ctx)
	if err != nil {
		return
	}
	defer consumer.Close()

	log.Println("Start DockedPlayerMoney ...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-consumer.Errors():
			log.Println("Error: DockedPlayerMoney failed: ", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "buy" {
				_ = h.playerUsecase.UpsertOffset(ctx, msg.Offset+1)

				req := new(player.CreatePlayerTransactionReq)

				if err = queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.playerUsecase.DockedPlayerMoneyRes(ctx, req)

				log.Printf("DockedPlayerMoney | Topic(%s)| Offset(%d) Message(%s) \n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigChan:
			log.Println("Stop DockedPlayerMoney...")
			return
		}
	}
}

func (h *playerQueueHandler) RollbackPlayerTransaction() {
	ctx := context.Background()

	consumer, err := h.playerConsumer(ctx)
	if err != nil {
		return
	}
	defer consumer.Close()

	log.Println("Start RollbackPlayerTransaction...")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-consumer.Errors():
			log.Printf("Error: RollbackPlayerTransaction: %v \n", err)
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) != "rollback_buy" {
				if err = h.playerUsecase.UpsertOffset(ctx, msg.Offset+1); err != nil {
					continue
				}

				req := new(player.RollbackPlayerTransactionReq)
				if err = queue.DecodeMessage(req, msg.Value); err != nil {
					log.Printf("Error: DecodeMessage: %v \n", err)
					continue
				}

				if err = h.playerUsecase.RollbackPlayerTransaction(ctx, req); err != nil {
					log.Printf("Error: RollbackPlayerTransaction: %v \n", err)
					continue
				}

				log.Printf("RollbackPlayerTransaction: Topic(%s) | Offset(%d) | Message(%s)  \n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigCh:
			log.Println("Stop RollbackPlayerTransaction...")
		}
	}
}
