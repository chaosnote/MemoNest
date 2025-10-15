package model

type CategoryModel struct {
	Categories []Category
}

type Category struct {
	RowID    int
	NodeID   string
	ParentID string
	PathName string
	LftIdx   int
	RftIdx   int
}

// 定義一個新的結構體來表示巢狀的分類樹
type CategoryNode struct {
	Category
	Children []*CategoryNode // 關鍵變更：使用指標切片
	Path     string          // 新增路徑欄位
}
