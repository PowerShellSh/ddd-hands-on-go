package book

import "fmt"

// Title は書籍のタイトルを表す値オブジェクトです。
type Title struct {
	value string
}

// NewTitle は新しいTitleを生成します。
// 値が空の場合はエラーを返します。
func NewTitle(value string) (*Title, error) {
	if value == "" {
		return nil, fmt.Errorf("Titleの生成に失敗しました: タイトルは必須です")
	}
	return &Title{value: value}, nil
}

// Value はTitleの値を返します。
func (t *Title) Value() string {
	return t.value
}
