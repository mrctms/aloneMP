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

type ProgressView struct {
	*cview.ProgressBar
}

func NewProgressView() *ProgressView {
	progressView := new(ProgressView)
	progressView.ProgressBar = cview.NewProgressBar()
	progressView.ProgressBar.FilledColor = tcell.ColorGreen
	progressView.ProgressBar.SetTitle("00:00:00/00:00:00")
	progressView.ProgressBar.SetBorder(true)
	progressView.ProgressBar.SetMax(100)
	return progressView
}

func (p *ProgressView) SetProgressTitle(t string) {
	p.ProgressBar.SetTitle(t)
}

func (p *ProgressView) UpdateProgress(progress int) {
	p.ProgressBar.SetProgress(progress)
}
