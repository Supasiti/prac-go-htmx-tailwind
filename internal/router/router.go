package router

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/supasiti/prac-go-htmx-tailwind/internal/model"
	"github.com/supasiti/prac-go-htmx-tailwind/templates/components"
	"github.com/supasiti/prac-go-htmx-tailwind/templates/page"
)

var (
	contacts = map[int]*model.Contact{
		1: {
			ContactID: 1,
			Name:      "Roman Gurevitch",
			Email:     "r.gurevitch@email.com",
		},
		2: {
			ContactID: 2,
			Name:      "Adrian Lai",
			Email:     "a.lai@email.com",
		},
		3: {
			ContactID: 3,
			Name:      "Jason Ngo",
			Email:     "j.ngo@email.com",
		},
	}
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	csh := NewContactsHandler()
	mux.Handle("/contact", csh)

	ch := NewContactHandler()
	mux.Handle("/contact/", ch)

	return mux
}

type contactsHandler struct{}

func NewContactsHandler() *contactsHandler {
	return &contactsHandler{}
}

func (h *contactsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.GET(w, r)
		return
	}
	http.NotFound(w, r)
}

func (h *contactsHandler) GET(w http.ResponseWriter, r *http.Request) {
	t := components.ContactTable(ToArray(contacts))
	page.Page(t).Render(r.Context(), w)
}

func ToArray(cs map[int]*model.Contact) []*model.Contact {
	result := []*model.Contact{}

	for _, c := range cs {
		result = append(result, c)
	}

	return result
}

type contactHandler struct{}

func NewContactHandler() *contactHandler {
	return &contactHandler{}
}

func (h *contactHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contactID, err := getContactIdFromURL(r)
	if err != nil {
		http.Error(w, "Invalid contactID", http.StatusBadRequest)
		return
	}
	slog.Info("GET", "contactID", contactID)

	if r.Method == http.MethodGet {
		h.GET(w, r, contactID)
		return
	}
	http.NotFound(w, r)
}

func (h *contactHandler) GET(w http.ResponseWriter, r *http.Request, contactID int) {
	contact, ok := contacts[contactID]
	if !ok {
		http.Error(w, "Invalid contactID", http.StatusBadRequest)
		return
	}

	action := r.URL.Query().Get("action")
	if action == "edit" {
		components.ContactForm(contact).Render(r.Context(), w)
		return
	}

	components.ContactRow(contact).Render(r.Context(), w)
}

func getContactIdFromURL(r *http.Request) (int, error) {
	parts := strings.Split(r.URL.Path, "/")
	slog.Info("getContactIDFromURL", "parts", parts)

	// first entry is always an empty string
	if len(parts) < 3 {
		return 0, errors.New("missing id")
	}

	contactID := parts[2]
	slog.Info("contactID", "id", contactID)
	return strconv.Atoi(contactID)
}
