// Package layout re-exports the TUIOS BSP tiling layout engine for external use.
package layout

import (
	il "github.com/Gaurav-Gosain/tuios/internal/layout"
)

// BSP tree types.
type (
	BSPTree   = il.BSPTree
	TileNode  = il.TileNode
	Rect      = il.Rect
	SplitType = il.SplitType

	AutoScheme      = il.AutoScheme
	PreselectionDir = il.PreselectionDir

	SerializedNode    = il.SerializedNode
	SerializedBSPTree = il.SerializedBSPTree

	TileLayout = il.TileLayout
)

// SplitType constants.
const (
	SplitNone       = il.SplitNone
	SplitVertical   = il.SplitVertical
	SplitHorizontal = il.SplitHorizontal
)

// AutoScheme constants.
const (
	SchemeLongestSide = il.SchemeLongestSide
	SchemeAlternate   = il.SchemeAlternate
	SchemeSpiral      = il.SchemeSpiral
)

// PreselectionDir constants.
const (
	PreselectionNone  = il.PreselectionNone
	PreselectionLeft  = il.PreselectionLeft
	PreselectionRight = il.PreselectionRight
	PreselectionUp    = il.PreselectionUp
	PreselectionDown  = il.PreselectionDown
)

// Constructors.
var (
	NewBSPTree           = il.NewBSPTree
	NewLeafNode          = il.NewLeafNode
	NewInternalNode      = il.NewInternalNode
	BuildTreeFromWindows = il.BuildTreeFromWindows
	ParseAutoScheme      = il.ParseAutoScheme
	CalculateTilingLayout = il.CalculateTilingLayout
)
