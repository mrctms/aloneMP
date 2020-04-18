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
