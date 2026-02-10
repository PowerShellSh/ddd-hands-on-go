package book

import "fmt"

// BookId は書籍のIDを表す値オブジェクトです。
type BookId struct {
	value string
}

// NewBookId は新しいBookIdを生成します。
// 値が空の場合はエラーを返します。
func NewBookId(value string) (*BookId, error) {
	if value == "" {
		return nil, fmt.Errorf("BookIdの生成に失敗しました: ISBNは必須です")
	}
	// TODO: 必要であればISBNのフォーマット検証を追加する
	return &BookId{value: value}, nil
}

// Value はBookIdの値を返します。
func (id *BookId) Value() string {
	return id.value
}

func (id *BookId) Equals(other *BookId) bool {
	if other == nil {
		return false
	}
	return id.value == other.value
}
