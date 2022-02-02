package server_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"pet-store-server/internal/petstore"
	"pet-store-server/internal/server"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	require := require.New(t)
	ctrl := gomock.NewController(t)
	petstoreMock := petstore.NewMockInterface(ctrl)
	service := server.NewService(petstoreMock)

	go service.Start(":12345")

	t.Log("Successful pets list")
	{
		petstoreMock.EXPECT().List().
			Return([]*petstore.Pet{
				{
					ID:    1,
					Type:  petstore.Fish,
					Price: 0.5,
				},
				{
					ID:    2,
					Type:  petstore.Bird,
					Price: 2.5,
				},
				{
					ID:    3,
					Type:  petstore.Cat,
					Price: 12.5,
				},
			}, nil)

		resp, err := http.Get("http://localhost:12345/pets")
		require.NoError(err)
		require.Equal(http.StatusOK, resp.StatusCode)

		data, err := ioutil.ReadAll(resp.Body)
		require.NoError(err)
		require.Equal(
			`[{"id":1,"type":"fish","price":0.5},{"id":2,"type":"bird","price":2.5},{"id":3,"type":"cat","price":12.5}]`,
			string(data),
		)
	}

	t.Log("Failed pets list")
	{
		petstoreMock.EXPECT().List().
			Return(nil, errors.New("some error"))

		resp, err := http.Get("http://localhost:12345/pets")
		require.NoError(err)
		require.Equal(http.StatusInternalServerError, resp.StatusCode)

		data, err := ioutil.ReadAll(resp.Body)
		require.NoError(err)
		require.Equal(
			`some error`,
			string(data),
		)
	}

	t.Log("Successful get pet")
	{
		petstoreMock.EXPECT().Get(uint64(12345)).
			Return(&petstore.Pet{
				ID:    12345,
				Type:  petstore.Gecko,
				Price: 129.99,
			}, nil)

		resp, err := http.Get("http://localhost:12345/pets/12345")
		require.NoError(err)
		require.Equal(http.StatusOK, resp.StatusCode)

		data, err := ioutil.ReadAll(resp.Body)
		require.NoError(err)
		require.Equal(
			`{"id":12345,"type":"gecko","price":129.99}`,
			string(data),
		)
	}

	t.Log("Failed get pet")
	{
		petstoreMock.EXPECT().Get(uint64(12345)).
			Return(nil, errors.New("some error"))

		resp, err := http.Get("http://localhost:12345/pets/12345")
		require.NoError(err)
		require.Equal(http.StatusInternalServerError, resp.StatusCode)

		data, err := ioutil.ReadAll(resp.Body)
		require.NoError(err)
		require.Equal(
			`some error`,
			string(data),
		)
	}

	t.Log("Invalid ID in get pet")
	{
		resp, err := http.Get("http://localhost:12345/pets/invalid-id")
		require.NoError(err)
		require.Equal(http.StatusNotFound, resp.StatusCode)
	}

	t.Log("Successful add pet")
	{
		petstoreMock.EXPECT().Add(petstore.Dog, float32(299.05)).
			Return(nil)

		resp, err := http.Post(
			"http://localhost:12345/pets",
			"application/json",
			bytes.NewReader([]byte(`{"type":"dog","price":299.05}`)),
		)
		require.NoError(err)
		require.Equal(http.StatusOK, resp.StatusCode)

		data, err := ioutil.ReadAll(resp.Body)
		require.NoError(err)
		require.Equal(
			``,
			string(data),
		)
	}

	t.Log("Failed add pet")
	{
		petstoreMock.EXPECT().Add(petstore.Dog, float32(299.05)).
			Return(errors.New("some error"))

		resp, err := http.Post(
			"http://localhost:12345/pets",
			"application/json",
			bytes.NewReader([]byte(`{"type":"dog","price":299.05}`)),
		)
		require.NoError(err)
		require.Equal(http.StatusInternalServerError, resp.StatusCode)

		data, err := ioutil.ReadAll(resp.Body)
		require.NoError(err)
		require.Equal(
			`some error`,
			string(data),
		)
	}

	t.Log("Inappropriate pet type in add pet")
	{
		resp, err := http.Post(
			"http://localhost:12345/pets",
			"application/json",
			bytes.NewReader([]byte(`{"type":"lizard","price":299.05}`)),
		)
		require.NoError(err)
		require.Equal(http.StatusBadRequest, resp.StatusCode)

		data, err := ioutil.ReadAll(resp.Body)
		require.NoError(err)
		require.Equal(
			`inappropriate pet type lizard`,
			string(data),
		)
	}

	t.Log("Inappropriate pet price in add pet")
	{
		resp, err := http.Post(
			"http://localhost:12345/pets",
			"application/json",
			bytes.NewReader([]byte(`{"type":"fish","price":-299.05}`)),
		)
		require.NoError(err)
		require.Equal(http.StatusBadRequest, resp.StatusCode)

		data, err := ioutil.ReadAll(resp.Body)
		require.NoError(err)
		require.Equal(
			`price cannot be below zero`,
			string(data),
		)
	}
}
