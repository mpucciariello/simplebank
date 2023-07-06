package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/micaelapucciariello/simplebank/token"
	"github.com/micaelapucciariello/simplebank/utils"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/micaelapucciariello/simplebank/db/mock"
	db "github.com/micaelapucciariello/simplebank/db/sqlc"
)

const _amount = 112

var (
	user1, _   = randomUser()
	user2, _   = randomUser()
	userARS, _ = randomUser()
	account1   = randomAccount(user1.Username)
	account2   = randomAccount(user2.Username)
	accountARS = randomAccount(userARS.Username)
)

func TestCreateTransferAPI(t *testing.T) {
	accountARS.Currency = utils.ARS

	transfer := db.TransferTxResult{
		Transfer: db.Transfer{
			ID:            utils.RandomInt(1, 1000),
			FromAccountID: account1.ID,
			ToAccountID:   account2.ID,
			Amount:        _amount,
		},
		FromAccountID: account1,
		ToAccountID:   account2,
		FromEntry: db.Entry{
			ID:        utils.RandomInt(1, 1000),
			Amount:    -_amount,
			AccountID: account1.ID,
		},
		ToEntry: db.Entry{
			ID:        utils.RandomInt(1, 1000),
			Amount:    _amount,
			AccountID: account2.ID,
		},
	}

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "happy path create transfer",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          _amount,
				"currency":        utils.USD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, _authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        _amount,
				}
				store.EXPECT().GetAccount(gomock.Any(), account1.ID).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), account2.ID).Times(1).Return(account2, nil)
				store.EXPECT().TransferTx(gomock.Any(), arg).Times(1).
					Return(transfer, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusOK, recorder.Code)
				validateResponseTransfer(t, recorder.Body, transfer)
			},
		},
		{
			name: "internal server error",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          _amount,
				"currency":        utils.USD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, _authorizationTypeBearer, user1.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account1.ID).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), account2.ID).Times(1).Return(account2, nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(1).
					Return(db.TransferTxResult{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "invalid username",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          _amount,
				"currency":        utils.USD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, _authorizationTypeBearer, "unauthorized token", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account1.ID).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), account2.ID).Times(1).Return(account2, nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(1).
					Return(db.TransferTxResult{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "no authorization",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          _amount,
				"currency":        utils.USD,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account1.ID).Times(0).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), account2.ID).Times(0).Return(account2, nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0).
					Return(db.TransferTxResult{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "error: mismatched currency",
			body: gin.H{
				"from_account_id": accountARS.ID,
				"to_account_id":   account1.ID,
				"amount":          _amount,
				"currency":        utils.ARS,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, _authorizationTypeBearer, userARS.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.TransferTxParams{
					FromAccountID: accountARS.ID,
					ToAccountID:   account1.ID,
					Amount:        _amount,
				}

				store.EXPECT().GetAccount(gomock.Any(), accountARS.ID).Times(1).Return(accountARS, nil)
				store.EXPECT().GetAccount(gomock.Any(), account1.ID).Times(1).Return(account1, nil)
				store.EXPECT().TransferTx(gomock.Any(), arg).Times(0).
					Return(db.TransferTxResult{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check response
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			url := fmt.Sprintf("/transfers")

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.token)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func validateResponseTransfer(t *testing.T, body *bytes.Buffer, trxr db.TransferTxResult) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var rspTransfer db.TransferTxResult
	err = json.Unmarshal(data, &rspTransfer)
	require.NoError(t, err)
	require.Equal(t, trxr, rspTransfer)
}
