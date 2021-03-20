package cli

import (
	"util"

	"github.com/gdamore/tcell"
	"gitlab.com/tslocum/cview"
)

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

func (t *TracksList) AddTrackList(trackList util.TrackListMessage) {
	root := cview.NewTreeNode(trackList.Root).SetColor(tcell.ColorRed)
	t.TreeView.SetRoot(root).SetCurrentNode(root)
	for _, track := range trackList.Content {
		t.popolateTreeView(root, track)
	}

}

func (t *TracksList) popolateTreeView(targetNode *cview.TreeNode, rootInfo util.RootInfo) {
	node := cview.NewTreeNode(rootInfo.Name).SetReference(rootInfo.Path).
		SetSelectable(true).SetColor(tcell.ColorBlue).Collapse()
	if rootInfo.IsDir {
		node.SetColor(tcell.ColorRed)
		for _, v := range rootInfo.Content {
			t.popolateTreeView(node, v)
		}
	}
	targetNode.AddChild(node)
}

func (t *TracksList) GetSelectedTrackName() string {
	selectedNode := t.TreeView.GetCurrentNode()
	selectedNodeReference := selectedNode.GetReference()
	path := selectedNodeReference.(string)
	return path
}

func (t *TracksList) NextTrack() {
	t.TreeView.Transform(cview.TransformNextItem)
	selectedTrack := t.GetSelectedTrackName()
	if selectedTrack == "" {
		t.TreeView.GetCurrentNode().Expand()
		t.NextTrack()
	}
}

func (t *TracksList) PreviousTrack() {
	t.TreeView.Transform(cview.TransformPreviousItem)
	selectedTrack := t.GetSelectedTrackName()
	if selectedTrack == "" {
		t.TreeView.GetCurrentNode().Expand()
		t.PreviousTrack()
	}
}
