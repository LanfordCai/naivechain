package utils

import "encoding/json"

const (
	MAX_BLOCK_SERIALIZED_SIZE = 1000000 // bytes = 1MB

	COINBASE_MATURITY = 2

	// TODO: 这是啥？
	// Accept blocks timestamped as being from the future, up to this amount
	MAX_FUTURE_BLOCK_TIME = 60 * 60 * 2

	// TODO: 这是啥？ --> 聪？
	BELUSHIS_PER_COIN = int(100e6)

	TOTAL_COINS = 21000000

	MAX_MONEY = BELUSHIS_PER_COIN * TOTAL_COINS

	// block 被挖出来的间隔
	TIME_BETWEEN_BLOCKS_IN_SECS_TARGET = 1 * 60

	// 难度调整时间区间
	DIFFICULTY_PERIOD_IN_SECS_TARGET = 60 * 60 * 60

	DIFFICULTY_PERIOD_IN_BLOCKS = DIFFICULTY_PERIOD_IN_SECS_TARGET / TIME_BETWEEN_BLOCKS_IN_SECS_TARGET

	// TODO:
	INITIAL_DIFFICULTY_BITS = 24

	// 区块奖励减半时间
	HALVE_SUBSIDY_AFTER_BLOCKS_NUM = 210000
)

func Serialize(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func Deserialize(serialized string, v *interface{}) error {
	err := json.Unmarshal([]byte(serialized), v)
	if err != nil {
		return err
	}
	return nil
}

