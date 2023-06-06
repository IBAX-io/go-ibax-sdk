/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/
package accounts

type Account struct {
	Address Address
	URL     URL
}

type URL struct {
	Scheme string // Protocol scheme to identify a capable account backend
	Path   string // Path for the backend to identify a unique entity
}

type Address string

// Wallet represents a software or hardware wallet that might contain one or more
// accounts (derived from the same seed).
type Wallet interface {
	// URL retrieves the canonical path under which this wallet is reachable. It is
	// used by upper layers to define a sorting order over all wallets from multiple
	// backends.
	// URL() URL

	// Status returns a textual status to aid the user in the current state of the
	// wallet. It also returns an error indicating any failure the wallet might have
	// encountered.
	Status() (string, error)

	// Open initializes access to a wallet instance. It is not meant to unlock or
	// decrypt account keys, rather simply to establish a connection to hardware
	// wallets and/or to access derivation seeds.
	//
	// The password parameter may or may not be used by the implementation of a
	// particular wallet instance. The reason there is no passwordless open method
	// is to strive towards a uniform wallet handling, oblivious to the different
	// backend providers.
	//
	// Please note, if you open a wallet, you must close it to release any allocated
	// resources (especially important when working with hardware wallets).
	Open(password string) error

	// Close releases any resources held by an open wallet instance.
	Close() error

	// Accounts retrieves the list of signing accounts the wallet is currently aware
	// of. For hierarchical deterministic wallets, the list will not be exhaustive,
	// rather only contain the accounts explicitly pinned during account derivation.
	Accounts() []Account

	// Contains returns whether an account is part of this particular wallet or not.
	Contains(account Account) bool

	// Derive attempts to explicitly derive a hierarchical deterministic account at
	// the specified derivation path. If requested, the derived account will be added
	// to the wallet's tracked account list.
	Derive(path DerivationPath, pin bool) (Account, error)

	// SelfDerive sets a base account derivation path from which the wallet attempts
	// to discover non zero accounts and automatically add them to list of tracked
	// accounts.
	//
	// Note, self derivation will increment the last component of the specified path
	// opposed to descending into a child path to allow discovering accounts starting
	// from non zero components.
	//
	// Some hardware wallets switched derivation paths through their evolution, so
	// this method supports providing multiple bases to discover old user accounts
	// too. Only the last base will be used to derive the next empty account.
	//
	// You can disable automatic account discovery by calling SelfDerive with a nil
	// chain state reader.
	// @todo [ ] SelfDerive
	// SelfDerive(bases []DerivationPath, chain ethereum.ChainStateReader)

	// SignData requests the wallet to sign the hash of the given data
	// It looks up the account specified either solely via its address contained within,
	// or optionally with the aid of any location metadata from the embedded URL field.
	//
	// If the wallet requires additional authentication to sign the request (e.g.
	// a password to decrypt the account, or a PIN code to verify the transaction),
	// an AuthNeededError instance will be returned, containing infos for the user
	// about which fields or actions are needed. The user may retry by providing
	// the needed details via SignDataWithPassword, or by other means (e.g. unlock
	// the account in a keystore).
	SignData(account Account, mimeType string, data []byte) ([]byte, error)

	// SignDataWithPassword is identical to SignData, but also takes a password
	// NOTE: there's a chance that an erroneous call might mistake the two strings, and
	// supply password in the mimetype field, or vice versa. Thus, an implementation
	// should never echo the mimetype or return the mimetype in the error-response
	SignDataWithPassword(account Account, password, mimeType string, data []byte) ([]byte, error)

	// SignText requests the wallet to sign the hash of a given piece of data, prefixed
	// by the Ethereum prefix scheme
	// It looks up the account specified either solely via its address contained within,
	// or optionally with the aid of any location metadata from the embedded URL field.
	//
	// If the wallet requires additional authentication to sign the request (e.g.
	// a password to decrypt the account, or a PIN code to verify the transaction),
	// an AuthNeededError instance will be returned, containing infos for the user
	// about which fields or actions are needed. The user may retry by providing
	// the needed details via SignTextWithPassword, or by other means (e.g. unlock
	// the account in a keystore).
	//
	// This method should return the signature in 'canonical' format, with v 0 or 1.
	SignText(account Account, text []byte) ([]byte, error)

	// SignTextWithPassword is identical to Signtext, but also takes a password
	SignTextWithPassword(account Account, password string, hash []byte) ([]byte, error)

	// SignTx requests the wallet to sign the given transaction.
	//
	// It looks up the account specified either solely via its address contained within,
	// or optionally with the aid of any location metadata from the embedded URL field.
	//
	// If the wallet requires additional authentication to sign the request (e.g.
	// a password to decrypt the account, or a PIN code to verify the transaction),
	// an AuthNeededError instance will be returned, containing infos for the user
	// about which fields or actions are needed. The user may retry by providing
	// the needed details via SignTxWithPassword, or by other means (e.g. unlock
	// the account in a keystore).
	//! SignTx(account Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)

	// SignTxWithPassword is identical to SignTx, but also takes a password
	//! SignTxWithPassword(account Account, password string, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)
}
