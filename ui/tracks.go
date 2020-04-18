/*
Copyright (C) MarckTomack <marcktomack@tutanota.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

package ui

import (
	"github.com/gdamore/tcell"
	"gitlab.com/tslocum/cview"
)

type TracksList struct {
	*cview.List
}

func NewTracksList() *TracksList {
	trackList := new(TracksList)
	trackList.List = cview.NewList()
	trackList.List.SetTitle("Track List")
	trackList.List.SetBorder(true)
	trackList.List.SetSelectedBackgroundColor(tcell.ColorBlue)
	return trackList
}

func (t *TracksList) AddItems(items []string) {
	for _, v := range items {
		t.List.AddItem(v, "", 0, nil)
	}
}

func (t *TracksList) GetSelectedItemText() string {
	index := t.List.GetCurrentItem()
	path, _ := t.List.GetItemText(index)
	return path
}

func (t *TracksList) NextItem() string {
	index := t.List.GetCurrentItem()
	if index >= t.List.GetItemCount()-1 {
		index = 0
	} else {
		index++
	}
	t.List.SetCurrentItem(index)
	return t.GetSelectedItemText()
}
