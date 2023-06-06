package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/micaelapucciariello/simplebank/db/mock"
	db "github.com/micaelapucciariello/simplebank/db/sqlc"
	"github.com/micaelapucciariello/simplebank/utils"
)

func TestGetAccountsAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:      "happy path",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusOK, recorder.Code)
				validateResponseAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "account not found",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "internal server error",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			recorder := httptest.NewRecorder()
			server := NewServer(store)

			url := fmt.Sprintf("/accounts/%d", account.ID)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			// check request
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomAccount() db.Account {
	account := db.Account{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomBalance(),
		Currency: utils.RandomCurrency(),
		ID:       utils.RandomInt(1, 1000),
	}

	return account
}

func validateResponseAccount(t *testing.T, body *bytes.Buffer, acc db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var rspAccount db.Account
	err = json.Unmarshal(data, &rspAccount)
	require.NoError(t, err)
	require.Equal(t, acc, rspAccount)
}
