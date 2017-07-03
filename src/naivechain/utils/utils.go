package utils

const MININING_DIFFICULTY = 6

func IsValidDifficulty(hash string) bool {
	zeroCount := 0
	for _, r := range hash {
		if r == '0'	{
			zeroCount++
		} else if zeroCount < MININING_DIFFICULTY {
			// 如果遇上非 0 值，就没必要继续算了，如果 0 的数目小于 难度，计算下一个nonce
			return false
		}

		// 不管是否遇上非 0 值，都判断下是否找到了正确的 nonce
		if zeroCount >= MININING_DIFFICULTY {
			return true
		}
	}
	return false
}

