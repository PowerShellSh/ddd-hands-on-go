package status

// StatusEnum は在庫ステータスの列挙型です。
type StatusEnum int

const (
	// InStock は在庫ありを表します。
	InStock StatusEnum = iota // 0
	// LowStock は残りわずかを表します。
	LowStock // 1
	// OutOfStock は在庫切れを表します。
	OutOfStock // 2
)

// String はStatusEnumの文字列表現を返します。
func (s StatusEnum) String() string {
	switch s {
	case InStock:
		return "IN_STOCK"
	case LowStock:
		return "LOW_STOCK"
	case OutOfStock:
		return "OUT_OF_STOCK"
	default:
		return "UNKNOWN"
	}
}

// Status は在庫ステータスを表す値オブジェクトです。
type Status struct {
	value StatusEnum
}

// NewStatus は新しいStatusを生成します。
func NewStatus(value StatusEnum) *Status {
	return &Status{value: value}
}

// Value はStatusの値を返します。
func (s *Status) Value() StatusEnum {
	return s.value
}

// IsOutOfStock は在庫切れかどうかを判定します。
func (s *Status) IsOutOfStock() bool {
	return s.value == OutOfStock
}

// ToStatusEnum は文字列からStatusEnumへ変換します。
func ToStatusEnum(s string) StatusEnum {
	switch s {
	case "IN_STOCK":
		return InStock
	case "LOW_STOCK":
		return LowStock
	case "OUT_OF_STOCK":
		return OutOfStock
	default:
		return OutOfStock // デフォルトは安全側に倒して在庫切れとする
	}
}
