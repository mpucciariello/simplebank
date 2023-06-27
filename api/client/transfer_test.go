package client

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
)

func TestCreateTransferAPI(t *testing.T) {
	account1 := randomAccount()
	account2 := randomAccount()
	transfer := db.Transfer{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        112,
	}

	testCases := []struct {
		name          string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "happy path create transfer",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateTransfer(gomock.Any(), gomock.Eq(db.CreateTransferParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        112,
				})).
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
				store.EXPECT().CreateTransfer(gomock.Any(), gomock.Eq(db.CreateTransferParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        112,
				})).
					Times(1).
					Return(db.Transfer{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("/transfers")

			body := fmt.Sprintf(`{"from_account_id": "%v", "to_account_id": "%v", "amount": %v}`, account1.ID, account2.ID, 112)
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

func validateResponseTransfer(t *testing.T, body *bytes.Buffer, acc db.Transfer) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var rspTransfer db.Transfer
	err = json.Unmarshal(data, &rspTransfer)
	require.NoError(t, err)
	require.Equal(t, acc, rspTransfer)
}
