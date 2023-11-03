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
	nextID = 4
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
	switch r.Method {

	case http.MethodGet:
		h.GET(w, r)
		return

	case http.MethodPost:
		h.POST(w, r)
		return

	default:
		http.NotFound(w, r)
	}
}

func (h *contactsHandler) GET(w http.ResponseWriter, r *http.Request) {
	t := components.ContactTable(toArray(contacts))
	page.Page(t).Render(r.Context(), w)
}

func (h *contactsHandler) POST(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	name := r.Form.Get("name")
	email := r.Form.Get("email")

	if len(name) == 0 || len(email) == 0 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	newContact := &model.Contact{
		ContactID: nextID,
		Name:      name,
		Email:     email,
	}
	contacts[nextID] = newContact

	// update the next available id
	nextID += 1
	components.ContactRow(newContact).Render(r.Context(), w)
	components.AddContactForm().Render(r.Context(), w)
}

func toArray(cs map[int]*model.Contact) []*model.Contact {
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
	contactID, err := parseContactID(r)
	if err != nil {
		http.Error(w, "Invalid contactID", http.StatusBadRequest)
		return
	}
	slog.Info("ServeHTTP", "contactID", contactID)

	contact, ok := contacts[contactID]
	if !ok {
		http.NotFound(w, r)
		return
	}

	switch r.Method {

	case http.MethodGet:
		h.GET(w, r, contact)
		return

	case http.MethodPatch:
		h.PATCH(w, r, contact)
		return

	case http.MethodDelete:
		h.DELETE(w, r, contact)
		return

	default:
		http.NotFound(w, r)
	}
}

func (h *contactHandler) GET(w http.ResponseWriter, r *http.Request, contact *model.Contact) {
	action := r.URL.Query().Get("action")
	if action == "edit" {
		components.ContactForm(contact).Render(r.Context(), w)
		return
	}

	components.ContactRow(contact).Render(r.Context(), w)
}

func (h *contactHandler) PATCH(w http.ResponseWriter, r *http.Request, contact *model.Contact) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// TODO validate form ?
	contact.Name = r.Form.Get("name")
	contact.Email = r.Form.Get("email")

	components.ContactRow(contact).Render(r.Context(), w)
}

func (h *contactHandler) DELETE(w http.ResponseWriter, r *http.Request, contact *model.Contact) {
	id := contact.ContactID
	delete(contacts, id)
}

func parseContactID(r *http.Request) (int, error) {
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
