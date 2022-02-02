package petstore_test

import (
	"pet-store-server/internal/petstore"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	require := require.New(t)

	client := petstore.NewClient("http://petstore-demo-endpoint.execute-api.com/petstore", 5*time.Second)
	require.NotNil(client)

	t.Log("Successful add")
	{
		err := client.Add(petstore.Dog, 20.00)
		require.NoError(err)
	}

	t.Log("Successfull get")
	{
		expected := &petstore.Pet{
			ID:    1,
			Type:  petstore.Dog,
			Price: 249.99,
		}
		pet, err := client.Get(1)
		require.NoError(err)
		require.Equal(expected, pet)
	}

	t.Log("Succesful list")
	{
		expected := []*petstore.Pet{
			{
				ID:    1,
				Type:  petstore.Dog,
				Price: 249.99,
			},
			{
				ID:    2,
				Type:  petstore.Cat,
				Price: 124.99,
			},
			{
				ID:    3,
				Type:  petstore.Fish,
				Price: 0.99,
			},
		}
		pets, err := client.List()
		require.NoError(err)
		require.Equal(expected, pets)
	}
}
