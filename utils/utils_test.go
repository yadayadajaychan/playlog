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

package utils_test

import (
	"testing"
	"github.com/yadayadajaychan/playlog/utils"
)

func TestScoreAndInternalLevelToDxRatingGen3(t *testing.T) {
	testCases := []struct {
		score         int
		internalLevel int
		dxRating      int
	}{
		{100_0470, 133, 287},
		{100_1379, 126, 272},
		{99_2970,  131, 270},
		{100_3532, 133, 288},
	}

	for i, tc := range testCases {
		rating := utils.ScoreAndInternalLevelToDxRatingGen3(tc.score, tc.internalLevel)
		if rating != tc.dxRating {
			t.Errorf("tc %d: expected %d, got %d", i, tc.dxRating, rating)
		}
	}
}
