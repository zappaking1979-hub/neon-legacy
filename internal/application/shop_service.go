package application

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/neonlegacy/server/internal/domain/item"
	"github.com/neonlegacy/server/internal/domain/player"
)

var (
	ErrNotEnoughCash    = fmt.Errorf("not enough cash")
	ErrItemNotFound     = fmt.Errorf("item not found")
	ErrNotEnoughItems   = fmt.Errorf("not enough items")
	ErrInvalidQuantity  = fmt.Errorf("invalid quantity")
)

type ShopService struct {
	itemRepo     item.ItemRepository
	playerItemRepo item.PlayerItemRepository
	playerRepo   player.Repository
}

func NewShopService(ir item.ItemRepository, pir item.PlayerItemRepository, pr player.Repository) *ShopService {
	return &ShopService{
		itemRepo:     ir,
		playerItemRepo: pir,
		playerRepo:   pr,
	}
}

func (s *ShopService) ListItems(ctx context.Context) ([]item.Item, error) {
	return s.itemRepo.List(ctx)
}

func (s *ShopService) GetItem(ctx context.Context, itemID int) (*item.Item, error) {
	return s.itemRepo.GetByID(ctx, itemID)
}

func (s *ShopService) ListInventory(ctx context.Context, playerID uuid.UUID) ([]item.PlayerItem, error) {
	return s.playerItemRepo.ListByPlayer(ctx, playerID)
}

func (s *ShopService) BuyItem(ctx context.Context, p *player.Player, itemID int, quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	it, err := s.itemRepo.GetByID(ctx, itemID)
	if err != nil {
		return fmt.Errorf("get item: %w", err)
	}
	if it == nil {
		return ErrItemNotFound
	}

	totalCost := it.BuyPrice * int64(quantity)
	if p.Cash < totalCost {
		return ErrNotEnoughCash
	}

	p.Cash -= totalCost

	if err := s.playerItemRepo.Add(ctx, p.ID, itemID, quantity); err != nil {
		return fmt.Errorf("add to inventory: %w", err)
	}

	if err := s.playerRepo.Update(ctx, p); err != nil {
		return fmt.Errorf("update player: %w", err)
	}

	return nil
}

func (s *ShopService) SellItem(ctx context.Context, p *player.Player, itemID int, quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	pi, err := s.playerItemRepo.GetByPlayerAndItem(ctx, p.ID, itemID)
	if err != nil {
		return fmt.Errorf("get player item: %w", err)
	}
	if pi == nil || pi.Quantity < quantity {
		return ErrNotEnoughItems
	}

	it, err := s.itemRepo.GetByID(ctx, itemID)
	if err != nil {
		return fmt.Errorf("get item: %w", err)
	}
	if it == nil {
		return ErrItemNotFound
	}

	totalGain := it.SellPrice * int64(quantity)
	p.Cash += totalGain

	if err := s.playerItemRepo.Remove(ctx, p.ID, itemID, quantity); err != nil {
		return fmt.Errorf("remove from inventory: %w", err)
	}

	if err := s.playerRepo.Update(ctx, p); err != nil {
		return fmt.Errorf("update player: %w", err)
	}

	return nil
}

type InventoryService struct {
	playerItemRepo item.PlayerItemRepository
	playerRepo     player.Repository
}

func NewInventoryService(pir item.PlayerItemRepository, pr player.Repository) *InventoryService {
	return &InventoryService{
		playerItemRepo: pir,
		playerRepo:     pr,
	}
}

func (s *InventoryService) ListInventory(ctx context.Context, playerID uuid.UUID) ([]item.PlayerItem, error) {
	return s.playerItemRepo.ListByPlayer(ctx, playerID)
}

func (s *InventoryService) UseItem(ctx context.Context, p *player.Player, itemID int) (string, error) {
	pi, err := s.playerItemRepo.GetByPlayerAndItem(ctx, p.ID, itemID)
	if err != nil {
		return "", fmt.Errorf("get player item: %w", err)
	}
	if pi == nil || pi.Quantity < 1 {
		return "", ErrNotEnoughItems
	}

	msg := pi.Use(p)

	if err := s.playerItemRepo.Remove(ctx, p.ID, itemID, 1); err != nil {
		return "", fmt.Errorf("remove item: %w", err)
	}

	if err := s.playerRepo.Update(ctx, p); err != nil {
		return "", fmt.Errorf("update player: %w", err)
	}

	return msg, nil
}