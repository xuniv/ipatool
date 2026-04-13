package core

import (
	"github.com/majd/ipatool/v2/pkg/appstore"
)

var (
	ErrAuthCodeRequired     = appstore.ErrAuthCodeRequired
	ErrPasswordTokenExpired = appstore.ErrPasswordTokenExpired
	ErrLicenseRequired      = appstore.ErrLicenseRequired
)

type (
	Error                   = appstore.Error
	Account                 = appstore.Account
	App                     = appstore.App
	Apps                    = appstore.Apps
	BagInput                = appstore.BagInput
	LoginInput              = appstore.LoginInput
	LookupInput             = appstore.LookupInput
	PurchaseInput           = appstore.PurchaseInput
	SearchInput             = appstore.SearchInput
	DownloadInput           = appstore.DownloadInput
	ReplicateSinfInput      = appstore.ReplicateSinfInput
	ListVersionsInput       = appstore.ListVersionsInput
	GetVersionMetadataInput = appstore.GetVersionMetadataInput
)

type Service interface {
	Bag(input BagInput) (appstore.BagOutput, error)
	Login(input LoginInput) (appstore.LoginOutput, error)
	AccountInfo() (appstore.AccountInfoOutput, error)
	Revoke() error
	Lookup(input LookupInput) (appstore.LookupOutput, error)
	Purchase(input PurchaseInput) error
	Search(input SearchInput) (appstore.SearchOutput, error)
	Download(input DownloadInput) (appstore.DownloadOutput, error)
	ReplicateSinf(input ReplicateSinfInput) error
	ListVersions(input ListVersionsInput) (appstore.ListVersionsOutput, error)
	GetVersionMetadata(input GetVersionMetadataInput) (appstore.GetVersionMetadataOutput, error)
}

type service struct {
	appStore appstore.AppStore
}

func New(appStore appstore.AppStore) Service {
	return service{appStore: appStore}
}

func (s service) Bag(input BagInput) (appstore.BagOutput, error) {
	return s.appStore.Bag(input)
}

func (s service) Login(input LoginInput) (appstore.LoginOutput, error) {
	return s.appStore.Login(input)
}

func (s service) AccountInfo() (appstore.AccountInfoOutput, error) {
	return s.appStore.AccountInfo()
}

func (s service) Revoke() error {
	return s.appStore.Revoke()
}

func (s service) Lookup(input LookupInput) (appstore.LookupOutput, error) {
	return s.appStore.Lookup(input)
}

func (s service) Purchase(input PurchaseInput) error {
	return s.appStore.Purchase(input)
}

func (s service) Search(input SearchInput) (appstore.SearchOutput, error) {
	return s.appStore.Search(input)
}

func (s service) Download(input DownloadInput) (appstore.DownloadOutput, error) {
	return s.appStore.Download(input)
}

func (s service) ReplicateSinf(input ReplicateSinfInput) error {
	return s.appStore.ReplicateSinf(input)
}

func (s service) ListVersions(input ListVersionsInput) (appstore.ListVersionsOutput, error) {
	return s.appStore.ListVersions(input)
}

func (s service) GetVersionMetadata(input GetVersionMetadataInput) (appstore.GetVersionMetadataOutput, error) {
	return s.appStore.GetVersionMetadata(input)
}
