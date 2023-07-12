package api

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/golang/mock/gomock"
	mockdb "github.com/micaelapucciariello/simplebank/db/mock"
	db "github.com/micaelapucciariello/simplebank/db/sqlc"
	"github.com/micaelapucciariello/simplebank/token"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRenewTokenAPI(t *testing.T) {
	user, _ := randomUser()
	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		account       db.Account
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "happy path renew token",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, _authorizationTypeBearer, user.Username, time.Minute)
			},
			account: account,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Eq(db.CreateAccountParams{
					Owner:    account.Owner,
					Balance:  0,
					Currency: account.Currency,
				})).
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
			name: "internal server error",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, _authorizationTypeBearer, user.Username, time.Minute)
			},
			account: account,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Eq(db.CreateAccountParams{
					Owner:    account.Owner,
					Balance:  0,
					Currency: account.Currency,
				})).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "unauthorized user",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, _authorizationTypeBearer, "unauthorized_user", time.Minute)
			},
			account: account,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Eq(db.CreateAccountParams{
					Owner:    account.Owner,
					Balance:  0,
					Currency: account.Currency,
				})).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "no authorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			account:   account,
			buildStubs: func(store *mockdb.MockStore) {
				// build stubs
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
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
			server := newTestServer(t, store)

			url := fmt.Sprintf("/accounts")

			body := fmt.Sprintf(`{"owner": "%v", "currency": "%v"}`, account.Owner, account.Currency)
			jsonBody := []byte(body)
			bodyReader := bytes.NewReader(jsonBody)

			request, err := http.NewRequest(http.MethodPost, url, bodyReader)
			// check request
			require.NoError(t, err)
			tc.setupAuth(t, request, server.token)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
