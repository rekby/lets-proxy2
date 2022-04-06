package acme_client_manager

import (
	"crypto/rsa"
	"encoding/json"

	"golang.org/x/crypto/acme"
	"golang.org/x/xerrors"
)

const stateFormatVersion = 1

type acmeManagerState struct {
	Accounts []acmeAccountState
	Version  int

	PrivateKeyDeprecated  *rsa.PrivateKey `json:"PrivateKey,omitempty"`
	AcmeAccountDeprecated *acme.Account   `json:"AcmeAccount,omitempty"`
}

func (state *acmeManagerState) Load(data []byte) (bool, error) {
	err := json.Unmarshal(data, state)
	if err != nil {
		return false, xerrors.Errorf("failed to unmarshal acme manager state: %w", err)
	}

	if migrated, err := state.migrate(); err == nil {
		return migrated, nil
	} else {
		return false, xerrors.Errorf("failed state migration: %w", err)
	}
}

func (state *acmeManagerState) check() error {
	if state == nil {
		return xerrors.Errorf("state is nil")
	}

	if state.PrivateKeyDeprecated != nil {
		return xerrors.Errorf("deprecated field not empty: PrivateKeyDeprecated")
	}
	if state.AcmeAccountDeprecated != nil {
		return xerrors.Errorf("deprecated field not empty: AcmeAccountDeprecated")
	}

	if state.Version != stateFormatVersion {
		return xerrors.Errorf("bad state version: %v", state.Version)
	}

	if len(state.Accounts) == 0 {
		return xerrors.Errorf("no accounts")
	}

	for index, account := range state.Accounts {
		if err := account.PrivateKey.Validate(); err != nil {
			return xerrors.Errorf("bad private key in account: %v", index)
		}

		if account.AcmeAccount.URI == "" {
			return xerrors.Errorf("empty acme account uri in account: %v", index)
		}
	}

	return nil
}

func (state *acmeManagerState) migrate() (bool, error) {
	hasMigrations := false
	for {
		switch state.Version {
		case 0:
			hasMigrations = true
			if state.AcmeAccountDeprecated != nil && state.PrivateKeyDeprecated != nil {
				state.Accounts = []acmeAccountState{{
					AcmeAccount: state.AcmeAccountDeprecated,
					PrivateKey:  state.PrivateKeyDeprecated,
				}}
			}
			state.Version++
		case 1:
			if stateFormatVersion != state.Version {
				return false, xerrors.Errorf("internal format version error: %v", stateFormatVersion)
			}
			return hasMigrations, nil
		default:
			return false, xerrors.Errorf("unexpected state version: %v", state.Version)
		}
	}
}

type acmeAccountState struct {
	PrivateKey  *rsa.PrivateKey `json:"PrivateKey"`
	AcmeAccount *acme.Account   `json:"AcmeAccount"`
}
