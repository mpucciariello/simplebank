package api

import (
	"database/sql"
	"fmt"
	"github.com/golang/mock/gomock"
	mockdb "github.com/micaelapucciariello/simplebank/db/mock"
	db "github.com/micaelapucciariello/simplebank/db/sqlc"
	"github.com/micaelapucciariello/simplebank/utils"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAccountsAPI(t *testing.T) {
	account := randomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// build stubs
	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(nil, account)

	server := New(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/%v", account.ID)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	// check request
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)
	// check response
	require.Equal(t, http.StatusOK, recorder.Code)
}

func randomAccount() db.Account {
	account := db.Account{
		Owner:     utils.RandomOwner(),
		Balance:   utils.RandomBalance(),
		Currency:  utils.RandomCurrency(),
		ID:        utils.RandomInt(0, 1000),
		CreatedAt: sql.NullTime{Time: time.Time{}},
	}

	return account
}
