package services

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/interfaces"

	"github.com/tomochain/tomox-sdk/types"
)

// TokenService struct with daos required, responsible for communicating with daos.
// TokenService functions are responsible for interacting with daos and implements business logics.
type TokenService struct {
	tokenDao interfaces.TokenDao
}

// NewTokenService returns a new instance of TokenService
func NewTokenService(tokenDao interfaces.TokenDao) *TokenService {
	return &TokenService{tokenDao}
}

// Create inserts a new token into the database
func (s *TokenService) Create(token *types.Token) error {
	t, err := s.tokenDao.GetByAddress(token.ContractAddress)
	if err != nil {
		logger.Error(err)
		return err
	}

	if t != nil {
		return ErrTokenExists
	}

	err = s.tokenDao.Create(token)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// GetByID fetches the detailed document of a token using its mongo ID
func (s *TokenService) GetByID(id bson.ObjectId) (*types.Token, error) {
	return s.tokenDao.GetByID(id)
}

// GetByAddress fetches the detailed document of a token using its contract address
func (s *TokenService) GetByAddress(addr common.Address) (*types.Token, error) {
	return s.tokenDao.GetByAddress(addr)
}

// GetAll fetches all the tokens from db
func (s *TokenService) GetAll() ([]types.Token, error) {
	return s.tokenDao.GetAll()
}

// GetQuote fetches all the quote tokens from db
func (s *TokenService) GetQuoteTokens() ([]types.Token, error) {
	return s.tokenDao.GetQuoteTokens()
}

// GetBase fetches all the quote tokens from db
func (s *TokenService) GetBaseTokens() ([]types.Token, error) {
	return s.tokenDao.GetBaseTokens()
}
