package domain

import (
	"context"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Song struct {
	ID         int64          `json:"id"`
	BundleID   int64          `json:"bundle_id"`
	CreatorID  int64          `json:"creator_id"`
	Title      string         `json:"title" gorm:"type:varchar(255);uniqueIndex:title_subtitle"`
	Subtitle   string         `json:"subtitle" gorm:"type:varchar(255);uniqueIndex:title_subtitle"`
	Key        string         `json:"key" gorm:"type:varchar(3);column:song_key"`
	Bpm        uint           `json:"bpm"`
	ChordSheet datatypes.JSON `json:"chord_sheet"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-"`
}

func (song Song) IsValidKey() bool {
	validKeys := [...]string{
		"C", "D", "E", "F", "G", "A", "B",
		"Cm", "Dm", "Em", "Fm", "Gm", "Am", "Bm",
		"C#", "D#", "F#", "G#", "A#",
		"C#m", "D#m", "F#m", "G#m", "A#m",
		"Db", "Eb", "Gb", "Ab", "Bb",
		"Dbm", "Ebm", "Gbm", "Abm", "Bbm",
	}

	for _, val := range validKeys {
		if song.Key == val {
			return true
		}
	}

	return false
}

type SongFilterOptions struct {
	IDs   []int64
	BIDs  []int64
	CIDs  []int64
	Title string
	Keys  []string
	Bpms  []uint
}

type SongService interface {
	Fetcher[Song]
	Fetch(ctx context.Context, options *SongFilterOptions) ([]Song, error)
	AuthSingleRemover[Song]
	AuthSingleStorer[Song]
	AuthSingleUpdater[Song]
}

type SongRepository interface {
	Getter[Song]
	Get(ctx context.Context, options *SongFilterOptions) ([]Song, error)
	Create(ctx context.Context, song *Song) error
	Delete(ctx context.Context, sid int64) error
	Update(ctx context.Context, song *Song) error
}
