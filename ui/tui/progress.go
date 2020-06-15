package tui

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
	progressView.ProgressBar.SetFilledColor(tcell.ColorGreen)
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
