package redismatching

import (
	"context"
	"fmt"
	"gameapp/entity"
	"gameapp/pkg/richerror"
	"gameapp/pkg/timestamp"
	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

// WaitingListPrefix TODO - add to config in usecase layer...
const WaitingListPrefix = "waitinglist"

func (d DB) AddToWaitingList(userID uint, category entity.Category) error {
	const op = richerror.Op("redismatching.AddToWaitingList")

	_, err := d.adapter.Client().ZAdd(context.Background(), fmt.Sprintf("%s:%s", WaitingListPrefix, category), redis.Z{
		Score:  float64(timestamp.Now()),
		Member: fmt.Sprintf("%d", userID),
	}).Result()
	if err != nil {
		return richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return nil
}

func (d DB) GetWaitingListByCategory(ctx context.Context, category entity.Category) ([]entity.WaitingMember, error) {
	const op = "redismatching.GetWaitingListByCategory"

	minimum := fmt.Sprintf("%d", timestamp.Add(-2*time.Hour))
	maximum := strconv.Itoa(int(timestamp.Now()))
	//panic("minimum and maximum")

	list, err := d.adapter.Client().ZRangeByScoreWithScores(ctx, getCategoryKey(category), &redis.ZRangeBy{
		Min:    minimum,
		Max:    maximum,
		Offset: 0,
		Count:  0,
	}).Result()
	if err != nil {
		return []entity.WaitingMember{}, richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	var result = make([]entity.WaitingMember, 0)
	for _, i := range list {
		UserID, _ := strconv.Atoi(i.Member.(string))

		result = append(result, entity.WaitingMember{
			UserID:    uint(UserID),
			Timestamp: int64(i.Score),
			Category:  category,
		})
	}

	return result, nil
}

func getCategoryKey(category entity.Category) string {
	return fmt.Sprintf("%s:%s", WaitingListPrefix, category)
}

func (d DB) RemoveUsersFromWaitingList(category entity.Category, userIDs []uint) {
	// TODO: add 5 to config
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	members := make([]any, 0)
	for _, u := range userIDs {
		members = append(members, strconv.Itoa(int(u)))
	}

	numberOfRemoveUser, err := d.adapter.Client().ZRem(ctx, getCategoryKey(category), members...).Result()
	if err != nil {
		log.Errorf("remove form waiting list : %v", err)
		// TODO: update metrics
	}

	log.Printf("%d items removed from %s", numberOfRemoveUser, getCategoryKey(category))
	// TODO: update metrics
}
