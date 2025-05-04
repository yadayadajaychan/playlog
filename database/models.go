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

// Package database handles the playlog and song database
package database

type Difficulty int
const (
	Basic Difficulty = iota
	Advanced
	Expert
	Master
	ReMaster
	Utage
)

type ComboStatus int
const (
	NoCombo ComboStatus = iota
	FullCombo
	FullComboPlus
	AllPerfect
	AllPerfectPlus
)

type SyncStatus int
const (
	NoSync SyncStatus = iota
	FullSync
	FullSyncPlus
	FullSyncDx
	FullSyncDxPlus
)

type SongInfo struct {
	SongId		int
	Name		string
	Artist		string
	Type		string
	Bpm		int
	Category	string
	Version		string
	Sort		string
	Charts		[]ChartInfo
}

type ChartInfo struct {
	Difficulty	Difficulty
	Level		int
	InternalLevel	int // multiplied by 10
	NotesDesigner	string
	MaxNotes	int
}

type PlaylogEntry struct {
	UserPlayDate	int64 // Unix timestamp

	SongId		int
	Difficulty	Difficulty
	Score		int
	DxScore		int
	ComboStatus	ComboStatus
	SyncStatus	SyncStatus
	IsClear		bool
	IsNewRecord	bool
	IsDxNewRecord	bool
	Track		int
	MatchingUsers	[]string

	MaxCombo	int
	TotalCombo	int
	MaxSync		int
	TotalSync	int

	FastCount	int
	LateCount	int
	BeforeRating	int
	AfterRating	int

	TapCriticalPerfect	int
        TapPerfect		int
	TapGreat		int
	TapGood			int
	TapMiss			int

	HoldCriticalPerfect	int
        HoldPerfect		int
	HoldGreat		int
	HoldGood		int
	HoldMiss		int

	SlideCriticalPerfect	int
        SlidePerfect		int
	SlideGreat		int
	SlideGood		int
	SlideMiss		int

	TouchCriticalPerfect	int
        TouchPerfect		int
	TouchGreat		int
	TouchGood		int
	TouchMiss		int

	BreakCriticalPerfect	int
        BreakPerfect		int
	BreakGreat		int
	BreakGood		int
	BreakMiss		int

	TotalCriticalPerfect	int
        TotalPerfect		int
	TotalGreat		int
	TotalGood		int
	TotalMiss		int
}
