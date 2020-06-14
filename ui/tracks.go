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
	"log"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell"
	"gitlab.com/tslocum/cview"
)

var knowExtension = [4]string{".mp3", ".wav", ".flac", ".ogg"}

type TracksList struct {
	*cview.TreeView
}

func NewTracksList() *TracksList {
	trackList := new(TracksList)
	trackList.TreeView = cview.NewTreeView()
	trackList.TreeView.SetTitle("Track List")
	trackList.TreeView.SetBorder(true)
	trackList.TreeView.SetSelectedFunc(func(node *cview.TreeNode) {
		node.SetExpanded(!node.IsExpanded())
	})
	return trackList
}

func (t *TracksList) AddItems(rootPath string) {
	root := cview.NewTreeNode(rootPath).SetColor(tcell.ColorRed)
	t.TreeView.SetRoot(root).SetCurrentNode(root)
	t.popolateTreeView(root, rootPath)
}

func (t *TracksList) popolateTreeView(targetNode *cview.TreeNode, rootPath string) {
	files := getKnowFiles(rootPath)

	for _, file := range files {
		node := cview.NewTreeNode(file.Name()).
			SetReference(filepath.Join(rootPath, file.Name())).
			SetSelectable(true).SetColor(tcell.ColorBlue).Collapse()
		if file.IsDir() {
			node.SetColor(tcell.ColorRed)
			t.popolateTreeView(node, filepath.Join(rootPath, file.Name()))
		}
		targetNode.AddChild(node)
	}
}

func (t *TracksList) GetSelectedTrackName() string {
	selectedNode := t.TreeView.GetCurrentNode()
	selectedNodeReference := selectedNode.GetReference()
	path := selectedNodeReference.(string)
	file, err := os.Stat(path)
	if err != nil {
		log.Fatalln(err)
	}
	if file.IsDir() {
		return ""
	} else {
		return path
	}
}

func (t *TracksList) NextTrack() {
	t.TreeView.Transform(cview.TransformNextItem)
	selectedTrack := t.GetSelectedTrackName()
	if selectedTrack == "" {
		t.TreeView.GetCurrentNode().Expand()
		t.NextTrack()
	}
}