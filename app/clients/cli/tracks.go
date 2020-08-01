package cli

import (
	"os"
	"path/filepath"

	"util"

	"github.com/gdamore/tcell"
	"gitlab.com/tslocum/cview"
)

type TracksList struct {
	*cview.TreeView
	allTrack []string
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
	files := util.GetKnowFilesInfo(rootPath)

	for _, file := range files {
		node := cview.NewTreeNode(file.Name()).
			SetReference(filepath.Join(rootPath, file.Name())).
			SetSelectable(true).SetColor(tcell.ColorBlue).Collapse()
		if file.IsDir() {
			node.SetColor(tcell.ColorRed)
			t.popolateTreeView(node, filepath.Join(rootPath, file.Name()))
		}
		t.internalLoadAllTracks(node)
		targetNode.AddChild(node)
	}
}

func (t *TracksList) internalLoadAllTracks(node *cview.TreeNode) {
	nodeReference := node.GetReference().(string)
	file, err := os.Stat(nodeReference)
	if err != nil {
		return
	}
	if !file.IsDir() {
		fileName := filepath.Base(nodeReference)
		t.allTrack = append(t.allTrack, fileName)
	}
}

func (t *TracksList) GetSelectedTrackName() string {
	selectedNode := t.TreeView.GetCurrentNode()
	selectedNodeReference := selectedNode.GetReference()
	path := selectedNodeReference.(string)
	file, err := os.Stat(path)
	if err != nil {
		return ""
	}
	if file.IsDir() {
		return ""
	} else {
		return path
	}
}

func (t *TracksList) GetAllTracks() []string {
	return t.allTrack
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
