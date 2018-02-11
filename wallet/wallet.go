package wallet

import (
	"github.com/tclchiam/oxidize-go/account"
	"github.com/tclchiam/oxidize-go/identity"
)

type Wallet interface {
	Identities() ([]*identity.Identity, error)
	NewIdentity() (*identity.Identity, error)
	Balance() ([]*account.Account, error)
}

type wallet struct {
	store *KeyStore
}

func NewWallet(store *KeyStore) Wallet {
	return &wallet{store: store}
}

func (w *wallet) Identities() ([]*identity.Identity, error) {
	return w.store.Identities()
}

func (w *wallet) NewIdentity() (*identity.Identity, error) {
	newIdentity := identity.RandomIdentity()
	err := w.store.SaveIdentity(newIdentity)
	if err != nil {
		return nil, err
	}

	return newIdentity, nil
}

func (w *wallet) Balance() ([]*account.Account, error) {
	return nil, nil
}
