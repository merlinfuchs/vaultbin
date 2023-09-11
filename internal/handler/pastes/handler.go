package pastes

import "github.com/merlinfuchs/vaultbin/internal/store"

type PastesHandler struct {
	store store.PastesStore
}

func New(store store.PastesStore) *PastesHandler {
	return &PastesHandler{
		store: store,
	}
}
