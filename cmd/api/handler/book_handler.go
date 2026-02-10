package handler

import (
	"ddd-hands-on-go/internal/application/book"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type BookHandler struct {
	registerBookService *book.RegisterBookApplicationService
	getBookService      *book.GetBookApplicationService
}

// NewBookHandler は新しいBookHandlerを生成します。
func NewBookHandler(
	registerBookService *book.RegisterBookApplicationService,
	getBookService *book.GetBookApplicationService,
) *BookHandler {
	return &BookHandler{
		registerBookService: registerBookService,
		getBookService:      getBookService,
	}
}

// RegisterBookRequest is the request body for registering a book.
// This struct is replaced by an anonymous struct in the Create method in the provided edit.
// Keeping it for now as the edit only shows a partial replacement.
// However, the instruction implies a full replacement of the RegisterBook method.
// Based on the provided `Code Edit` block, the `RegisterBookRequest` type is effectively removed
// and replaced by an anonymous struct within the `Create` method.
// I will remove the `RegisterBookRequest` type as it's no longer used by the new `Create` method.

// Create は書籍登録リクエストを処理します。
func (h *BookHandler) RegisterBook(w http.ResponseWriter, r *http.Request) { // Renamed from Create to RegisterBook to match original function name
	if r.Method != http.MethodPost {
		http.Error(w, "許可されていないメソッドです", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ISBN  string  `json:"isbn"`
		Title string  `json:"title"`
		Price float64 `json:"price"` // Renamed from PriceAmount to Price
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "リクエストボディの読み込みに失敗しました", http.StatusBadRequest)
		return
	}

	cmd := book.RegisterBookCommand{
		ISBN:        req.ISBN,
		Title:       req.Title,
		PriceAmount: req.Price, // Mapped to req.Price
	}

	if err := h.registerBookService.Execute(r.Context(), cmd); err != nil { // Changed from registerBookAppSvc to registerBookService
		// 本来はエラーの種類に応じてステータスコードを使い分けるべき
		log.Printf("書籍登録エラー: %v", err)
		http.Error(w, fmt.Sprintf("書籍の登録に失敗しました: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "書籍が正常に登録されました")
}

func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "許可されていないメソッドです", http.StatusMethodNotAllowed) // Localized error message
		return
	}

	isbn := r.PathValue("isbn")
	if isbn == "" {
		http.Error(w, "ISBN is required", http.StatusBadRequest)
		return
	}

	dto, err := h.getBookService.Execute(r.Context(), isbn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if dto == nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
