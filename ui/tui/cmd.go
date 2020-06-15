package tui

import "gitlab.com/tslocum/cview"

type CmdView struct {
	*cview.TextView
}

func NewCmdView() *CmdView {
	cmdView := new(CmdView)
	cmdView.TextView = cview.NewTextView().SetTextAlign(cview.AlignRight)
	cmdView.TextView.SetBorder(false)
	cmdView.TextView.SetDynamicColors(true)
	return cmdView
}

func (c *CmdView) SetText(text string) {
	c.TextView.SetText(text)
}
