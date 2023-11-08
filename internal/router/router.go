package router

import (
	"log/slog"
	"net/http"
	"strconv"

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

	csh := NewHandler()
	mux.Handle("/contact", csh)

	return mux
}

type handler struct{}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contactIDQuery := r.URL.Query().Get("contactID")
	if len(contactIDQuery) > 0 {
		h.serveOne(w, r, contactIDQuery)
		return
	}

	// No contactID
	switch r.Method {

	case http.MethodGet:
		slog.Info("GET /contact")
		h.GetAll(w, r)
		return

	case http.MethodPost:
		slog.Info("POST /contact")
		h.CreateOne(w, r)
		return

	default:
		http.NotFound(w, r)
	}
}

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	t := components.ContactTable(toArray(contacts))
	page.Page(t).Render(r.Context(), w)
}

func (h *handler) CreateOne(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	slog.Info("h.CreateOne", "form data", r.Form)

	name := r.Form.Get("name")
	email := r.Form.Get("email")

	if len(name) == 0 || len(email) == 0 {
		slog.Error("h.CreateOne", "error", "empty name or email")
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

func (h *handler) serveOne(w http.ResponseWriter, r *http.Request, q string) {
	contactID, err := strconv.Atoi(q)
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
		h.GetOne(w, r, contact)
		return

	case http.MethodPatch:
		h.UpdateOne(w, r, contact)
		return

	case http.MethodDelete:
		h.DeleteOne(w, r, contact)
		return

	default:
		http.NotFound(w, r)
	}
}

func (h *handler) GetOne(w http.ResponseWriter, r *http.Request, contact *model.Contact) {
	action := r.URL.Query().Get("action")
	if action == "edit" {
		components.ContactForm(contact).Render(r.Context(), w)
		return
	}

	components.ContactRow(contact).Render(r.Context(), w)
}

func (h *handler) UpdateOne(w http.ResponseWriter, r *http.Request, contact *model.Contact) {
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

func (h *handler) DeleteOne(w http.ResponseWriter, r *http.Request, contact *model.Contact) {
	id := contact.ContactID
	delete(contacts, id)
}
