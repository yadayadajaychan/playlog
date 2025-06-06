// Copyright (C) 2025 Ethan Cheng <ethan@nijika.org>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package utils

import (
	"math"
)

func ScoreToRank(score int) string {
	var rank string

	switch {
	case score >= 100_5000:
		rank = "SSS+"
	case score >= 100_0000:
		rank = "SSS"
	case score >= 99_5000:
		rank = "SS+"
	case score >= 99_0000:
		rank = "SS"
	case score >= 98_0000:
		rank = "S+"
	case score >= 97_0000:
		rank = "S"
	case score >= 94_0000:
		rank = "AAA"
	case score >= 90_0000:
		rank = "AA"
	case score >= 80_0000:
		rank = "A"
	case score >= 75_0000:
		rank = "BBB"
	case score >= 70_0000:
		rank = "BB"
	case score >= 60_0000:
		rank = "B"
	case score >= 50_0000:
		rank = "C"
	default:
		rank = "D"
	}

	return rank
}

// https://listed.to/@donmai/45107/exploring-the-algorithm-behind-maimai-dx-s-scoring-and-dx-rating-computation
func scoreToMultiplier(score int) float64 {
	var mult float64

	switch {
	case score >= 100_5000:
		mult = 22.4
	case score >= 100_4999:
		mult = 22.2
	case score >= 100_0000:
		mult = 21.6
	case score >= 99_9999:
		mult = 21.4
	case score >= 99_5000:
		mult = 21.1
	case score >= 99_0000:
		mult = 20.8
	case score >= 98_9999:
		mult = 20.6
	case score >= 98_0000:
		mult = 20.3
	case score >= 97_0000:
		mult = 20.0
	case score >= 96_9999:
		mult = 17.6
	case score >= 94_0000:
		mult = 16.8
	case score >= 90_0000:
		mult = 15.2
	case score >= 80_0000:
		mult = 13.6
	case score >= 79_9999:
		mult = 12.8
	case score >= 75_0000:
		mult = 12.0
	case score >= 70_0000:
		mult = 11.2
	case score >= 60_0000:
		mult = 9.6
	case score >= 50_0000:
		mult = 8.0
	case score >= 40_0000:
		mult = 6.4
	case score >= 30_0000:
		mult = 4.8
	case score >= 20_0000:
		mult = 3.2
	case score >= 10_0000:
		mult = 1.6
	default:
		mult = 0.0
	}

	return mult
}

// https://silentblue.remywiki.com/maimai_DX:Rating
func ScoreAndInternalLevelToDxRatingGen3(score int, internalLevel int) int {
	mult := scoreToMultiplier(score)
	sc := float64(score) / 1_000_000
	lvl := float64(internalLevel) / 10

	return int(math.Floor(mult * sc * lvl))
}
