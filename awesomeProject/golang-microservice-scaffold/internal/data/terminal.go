package data

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"example/internal/ent"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type TerminalData struct {
	db    *sql.Driver
	cache *redis.Client
}

// NewTerminalData create a TerminalData entity
func NewTerminalData(db *sql.Driver, cache *redis.Client) *TerminalData {
	return &TerminalData{db: db, cache: cache}
}

// GetTerminalByID if cache don't exist query db
func (d *TerminalData) GetTerminalByID(ctx context.Context, id int) (*ent.Terminal, error) {
	var terminal ent.Terminal

	// get from cache
	err := d.cache.Get(ctx, "terminal:"+string(id)).Scan(&terminal)
	if err == redis.Nil {
		log.Println("缓存未命中，尝试从数据库中获取终端信息")
		terminal, err := d.GetTerminalFromDB(ctx, id)
		if err != nil {
			return nil, err
		}
		// write in cache
		err = d.cache.Set(ctx, "terminal:"+string(id), terminal, 5*time.Minute).Err()
		if err != nil {
			log.Println("缓存写入失败: ", err)
		}
		return &terminal, nil
	} else if err != nil {
		return nil, err
	}
	return &terminal, nil
}

// GetTerminalFromDB
func (d *TerminalData) GetTerminalFromDB(ctx context.Context, id int) (entity.Terminal, error) {
	// query
	terminal, err := d.db.Terminal.Query().Where(terminal.ID(id)).Only(ctx)
	if err != nil {
		return entity.Terminal{}, err
	}
	return *terminal, nil
}

// UpdateTerminal
func (d *TerminalData) UpdateTerminal(ctx context.Context, id int, status string, timeout time.Duration) error {
	// update information of terminal
	if err := d.UpdateTerminalInDB(ctx, id, status, timeout); err != nil {
		return err
	}

	// update cache
	if err := d.cache.Del(ctx, "terminal:"+string(id)).Err(); err != nil {
		log.Println("缓存更新失败: ", err)
	}

	return nil
}

// UpdateTerminalInDB
func (d *TerminalData) UpdateTerminalInDB(ctx context.Context, id int, status string, timeout time.Duration) error {
	// update terminal
	_, err := d.db.Terminal.UpdateOneID(id).SetStatus(status).SetTimeout(int(timeout.Seconds())).Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

// SetTerminalTimeout
func (d *TerminalData) SetTerminalTimeout(ctx context.Context, id int, timeout time.Duration) error {
	//update terminal timeout
	if err := d.UpdateTerminalTimeoutInDB(ctx, id, timeout); err != nil {
		return err
	}

	// delete cache to make sure next time is the latest
	if err := d.cache.Del(ctx, "terminal:"+string(id)).Err(); err != nil {
		log.Println("缓存删除失败: ", err)
	}
	return nil
}

// UpdateTerminalTimeoutInDB
func (d *TerminalData) UpdateTerminalTimeoutInDB(ctx context.Context, id int, timeout time.Duration) error {

	_, err := d.db.Terminal.UpdateOneID(id).SetTimeout(int(timeout.Seconds())).Save(ctx)
	if err != nil {
		return err
	}
	return nil
}
