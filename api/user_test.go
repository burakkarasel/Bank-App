package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mockdb "github.com/burakkarasel/Bank-App/db/mock"
	db "github.com/burakkarasel/Bank-App/db/sqlc"
	"github.com/burakkarasel/Bank-App/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

// eqCreateUser struct implements gomock.Matcher interface
type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

// Matches implements gomock.Matcher interface
func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	// here we convert the given argument to CreateUserParams and check for error
	arg, ok := x.(db.CreateUserParams)

	if !ok {
		return false
	}

	// then we check if the password in eqCreateUserParamsMatcher matches with the hashedPassword in arg
	err := util.CheckPassword(e.password, arg.HashedPassword)

	if err != nil {
		return false
	}

	// if they matches we changed args hashedPassword to given hashed password
	e.arg.HashedPassword = arg.HashedPassword

	// and then we check if both password are strictly same
	return reflect.DeepEqual(e.arg, arg)
}

// String implements gomock.Matcher interface
func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

// EqCreateUserParams returns gomock.Matcher interface
func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

// TestCreateUserAPI tests CreateUser handler
func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"full_name": user.FullName,
				"email":     user.Email,
				"password":  password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requiredBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "Internal Error",
			body: gin.H{
				"username":  user.Username,
				"full_name": user.FullName,
				"email":     user.Email,
				"password":  password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Unique violation",
			body: gin.H{
				"username":  user.Username,
				"full_name": user.FullName,
				"email":     user.Email,
				"password":  password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "No username",
			body: gin.H{
				"full_name": user.FullName,
				"email":     user.Email,
				"password":  password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid username",
			body: gin.H{
				"username":  "abc#k",
				"full_name": user.FullName,
				"password":  password,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "No email",
			body: gin.H{
				"username":  user.Username,
				"full_name": user.FullName,
				"password":  password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid email",
			body: gin.H{
				"username":  user.Username,
				"full_name": user.FullName,
				"password":  password,
				"email":     "bck",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "No full name",
			body: gin.H{
				"username": user.Username,
				"password": password,
				"email":    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "No password",
			body: gin.H{
				"username":  user.Username,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Short Password",
			body: gin.H{
				"username":  user.Username,
				"full_name": user.FullName,
				"password":  "123",
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// create a new controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// create a new mockstore
			store := mockdb.NewMockStore(ctrl)
			tt.buildStubs(store)

			// create a new server with the mockstore
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/users"

			// marshal the body
			json, err := json.Marshal(tt.body)
			require.NoError(t, err)

			// create a request
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(json))
			require.NoError(t, err)

			// check for results
			server.router.ServeHTTP(recorder, req)
			tt.checkResponse(recorder)
		})
	}
}

func randomUser(t *testing.T) (db.User, string) {
	return db.User{
		Username: util.RandomOwner(),
		Email:    util.RandomEmail(),
		FullName: util.RandomString(12),
	}, "password"
}

func requiredBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	require.Equal(t, user, gotUser)
}
