package client

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-playground/locales/currency"
	"github.com/golang/mock/gomock"
	"github.com/micaelapucciariello/simplebank/utils"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/micaelapucciariello/simplebank/db/mock"
	db "github.com/micaelapucciariello/simplebank/db/sqlc"
)

const _amount = 112

func TestCreateTransferAPI(t *testing.T) {
	account1 := randomAccount()
	account2 := randomAccount()
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
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "happy path create transfer",
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        112,
				}
				store.EXPECT().TransferTx(gomock.Any(), arg).
					Times(1).
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
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(db.TransferTxParams{})).
					Times(1).
					Return(db.TransferTxResult{}, sql.ErrConnDone)
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
			server := newTestServer(t, store)

			url := fmt.Sprintf("/transfers")

			body := fmt.Sprintf(`{"from_account_id": "%v", "to_account_id": "%v", "amount": %v, "currency": %v"}`, account1.ID, account2.ID, _amount, currency.USD)
			jsonBody := []byte(body)
			bodyReader := bytes.NewReader(jsonBody)

			request, err := http.NewRequest(http.MethodPost, url, bodyReader)
			// check request
			require.NoError(t, err)

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
