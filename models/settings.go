package models

import (
	"log/slog"
	"strconv"
)

type Settings struct {
	Count  int
	Sort   string
	Ranked bool
}

func (s *Settings) SetRanked(c string) {
	s.Ranked = c != "false"
}
func (s *Settings) SetSort(c string) {
	switch c {
	case "1":
		s.Sort = "top"
	case "2":
		s.Sort = "recent"
	default:
		s.Sort = "top"
	}
}
func (s *Settings) SetCount(c string) {
	lCount, err := strconv.Atoi(c)
	if err != nil {
		s.Count = 100
	}
	s.Count = lCount
}
