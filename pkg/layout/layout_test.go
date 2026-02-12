package layout_test

import (
	"testing"

	"github.com/Gaurav-Gosain/tuios/pkg/layout"
)

func TestNewBSPTree(t *testing.T) {
	tree := layout.NewBSPTree()
	if tree == nil {
		t.Fatal("NewBSPTree returned nil")
	}
	if !tree.IsEmpty() {
		t.Error("New tree should be empty")
	}
}

func TestInsertAndLayout(t *testing.T) {
	tree := layout.NewBSPTree()
	bounds := layout.Rect{X: 0, Y: 0, W: 120, H: 40}

	tree.InsertWindow(1, 0, layout.SplitNone, 0.5, bounds)
	tree.InsertWindow(2, 1, layout.SplitVertical, 0.5, bounds)

	if tree.WindowCount() != 2 {
		t.Errorf("WindowCount = %d, want 2", tree.WindowCount())
	}

	rects := tree.ApplyLayout(bounds)
	if len(rects) != 2 {
		t.Errorf("ApplyLayout returned %d rects, want 2", len(rects))
	}
}

func TestSerializeDeserialize(t *testing.T) {
	tree := layout.NewBSPTree()
	bounds := layout.Rect{X: 0, Y: 0, W: 120, H: 40}
	tree.InsertWindow(1, 0, layout.SplitNone, 0.5, bounds)
	tree.InsertWindow(2, 1, layout.SplitVertical, 0.5, bounds)

	serialized := tree.Serialize()
	restored := serialized.Deserialize()

	if restored.WindowCount() != tree.WindowCount() {
		t.Errorf("Restored WindowCount = %d, want %d", restored.WindowCount(), tree.WindowCount())
	}
}
