package inventoryHandlers

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/korvised/ilog-shop/config"
	"github.com/korvised/ilog-shop/modules/inventory"
	"github.com/korvised/ilog-shop/modules/inventory/inventoryUsecases"
	"github.com/korvised/ilog-shop/pkg/queue"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type (
	InventoryQueueHandlerService interface {
		inventoryConsumer(c context.Context) (sarama.PartitionConsumer, error)
		AddPlayerItem()
		RemovePlayerItem()
		RollbackAddPlayerItem()
		RollbackRemovePlayerItem()
	}

	inventoryQueueHandler struct {
		cfg              *config.Config
		inventoryUsecase inventoryUsecases.InventoryUsecaseService
	}
)

func NewInventoryQueueHandler(cfg *config.Config, inventoryUsecase inventoryUsecases.InventoryUsecaseService) InventoryQueueHandlerService {
	return &inventoryQueueHandler{cfg, inventoryUsecase}
}

func (h *inventoryQueueHandler) inventoryConsumer(c context.Context) (sarama.PartitionConsumer, error) {
	worker, err := queue.ConnectConsumer([]string{h.cfg.Kafka.Url}, h.cfg.Kafka.ApiKey, h.cfg.Kafka.Secret)
	if err != nil {
		return nil, err
	}

	offset, err := h.inventoryUsecase.GetOffset(c)
	if err != nil {
		return nil, err
	}
	consumer, err := worker.ConsumePartition("inventory", 0, offset)
	if err != nil {
		log.Printf("Error inventoryConsumer 1: %v \n", err)
		consumer, err = worker.ConsumePartition("player", 0, offset)
		if err != nil {
			log.Printf("Error inventoryConsumer 2 : %v \n", err)
			return nil, err
		}
	}

	return consumer, nil
}

func (h *inventoryQueueHandler) AddPlayerItem() {
	ctx := context.Background()

	consumer, err := h.inventoryConsumer(ctx)
	if err != nil {
		return
	}
	defer consumer.Close()

	log.Println("Start AddPlayerItem Consumer ...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-consumer.Errors():
			log.Println("Error: AddPlayerItem failed: ", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "buy" {
				_ = h.inventoryUsecase.UpsertOffset(ctx, msg.Offset+1)

				req := new(inventory.UpdateInventoryReq)

				if err = queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.inventoryUsecase.AddPlayerItemRes(ctx, req)

				log.Printf("AddPlayerItem | Topic(%s)| Offset(%d) Message(%s) \n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigChan:
			log.Println("Stop AddPlayerItem...")
			return
		}
	}
}

func (h *inventoryQueueHandler) RemovePlayerItem() {
	ctx := context.Background()

	consumer, err := h.inventoryConsumer(ctx)
	if err != nil {
		return
	}
	defer consumer.Close()

	log.Println("Start RemovePlayerItem Consumer ...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-consumer.Errors():
			log.Println("Error: RemovePlayerItem failed: ", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "sell" {
				_ = h.inventoryUsecase.UpsertOffset(ctx, msg.Offset+1)

				req := new(inventory.UpdateInventoryReq)

				if err = queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.inventoryUsecase.RemovePlayerItemRes(ctx, req)

				log.Printf("RemovePlayerItem | Topic(%s)| Offset(%d) Message(%s) \n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigChan:
			log.Println("Stop RemovePlayerItem...")
			return
		}
	}
}

func (h *inventoryQueueHandler) RollbackAddPlayerItem() {
	ctx := context.Background()

	consumer, err := h.inventoryConsumer(ctx)
	if err != nil {
		return
	}
	defer consumer.Close()

	log.Println("Start RollbackAddPlayerItem Consumer ...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-consumer.Errors():
			log.Println("Error: RollbackAddPlayerItem failed: ", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "rollback_buy" {
				_ = h.inventoryUsecase.UpsertOffset(ctx, msg.Offset+1)

				req := new(inventory.RollbackPlayerInventoryReq)

				if err = queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.inventoryUsecase.RollbackAddPlayerItem(ctx, req)

				log.Printf("RollbackAddPlayerItem | Topic(%s)| Offset(%d) Message(%s) \n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigChan:
			log.Println("Stop RollbackAddPlayerItem...")
			return
		}
	}
}

func (h *inventoryQueueHandler) RollbackRemovePlayerItem() {
	ctx := context.Background()

	consumer, err := h.inventoryConsumer(ctx)
	if err != nil {
		return
	}
	defer consumer.Close()

	log.Println("Start RollbackRemovePlayerItem Consumer ...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-consumer.Errors():
			log.Println("Error: RollbackRemovePlayerItem failed: ", err.Error())
			continue
		case msg := <-consumer.Messages():
			if string(msg.Key) == "rollback_sell" {
				_ = h.inventoryUsecase.UpsertOffset(ctx, msg.Offset+1)

				req := new(inventory.RollbackPlayerInventoryReq)

				if err = queue.DecodeMessage(req, msg.Value); err != nil {
					continue
				}

				h.inventoryUsecase.RollbackRemovePlayerItem(ctx, req)

				log.Printf("RollbackRemovePlayerItem | Topic(%s)| Offset(%d) Message(%s) \n", msg.Topic, msg.Offset, string(msg.Value))
			}
		case <-sigChan:
			log.Println("Stop RollbackRemovePlayerItem...")
			return
		}
	}
}
