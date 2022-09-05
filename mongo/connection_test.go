package mongo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Connection(t *testing.T) {
	mongoURI := os.Getenv("local")
	_ = os.Getenv("user")
	_ = os.Getenv("password")

	t.Run("connection success", func(t *testing.T) {
		var mongoSecondary *Client
		var mongoPrimary *Client

		mongoSecondary = Connect(mongoURI, "sample_analytics", SECONDARY_PREFERRED, 1)
		mongoPrimary = Connect(mongoURI, "sample_analytics", SECONDARY_PREFERRED, 1)
		defer mongoPrimary.Disconnect()
		defer mongoSecondary.Disconnect()
		assert.NotNil(t, mongoPrimary.Client)
		assert.NotNil(t, mongoSecondary.Client)
	})

	t.Run("connection failed - host not found", func(t *testing.T) {
		var mongoSecondary *Client
		defer func() {
			r := recover()
			assert.NotNil(t, r)
		}()
		mongoSecondary = Connect("", "sample_analytics", SECONDARY_PREFERRED, 1)
		assert.Nil(t, mongoSecondary.Client)
	})

	t.Run("connection failed - password", func(t *testing.T) {
		var mongoSecondary *Client
		defer func() {
			r := recover()
			assert.NotNil(t, r)
		}()
		mongoSecondary = Connect(mongoURI, "sample_analytics", SECONDARY_PREFERRED, 1)
		assert.Nil(t, mongoSecondary.Client)
	})

	t.Run("connection failed - auth source", func(t *testing.T) {
		var mongoSecondary *Client
		defer func() {
			r := recover()
			assert.NotNil(t, r)
		}()

		mongoSecondary = Connect(mongoURI, "", SECONDARY_PREFERRED, 1)
		assert.Nil(t, mongoSecondary.Client)
	})

	t.Run("connection failed - user", func(t *testing.T) {
		var mongoSecondary *Client
		defer func() {
			r := recover()
			assert.NotNil(t, r)
		}()

		mongoSecondary = Connect(mongoURI, "sample_analytics", SECONDARY_PREFERRED, 1)
		assert.Nil(t, mongoSecondary.Client)
	})
}

func Test_GetDatabase(t *testing.T) {
	mongoURI := os.Getenv("local")

	var mongoSecondary *Client
	var mongoPrimary *Client
	mongoSecondary = Connect(mongoURI, "sample_analytics", SECONDARY_PREFERRED, 1)
	mongoPrimary = Connect(mongoURI, "sample_analytics", PRIMARY, 1)
	defer mongoPrimary.Disconnect()
	defer mongoSecondary.Disconnect()
	assert.NotNil(t, mongoPrimary.Client)
	assert.NotNil(t, mongoSecondary.Client)

	sec := mongoSecondary.GetDatabase("sample_analytics")
	prim := mongoPrimary.GetDatabase("sample_analytics")

	assert.NotNil(t, sec)
	assert.NotNil(t, prim)
}

func Test_GetCollection(t *testing.T) {
	mongoURI := os.Getenv("local")

	var mongoSecondary *Client
	var mongoPrimary *Client
	mongoSecondary = Connect(mongoURI, "sample_analytics", SECONDARY_PREFERRED, 1)
	mongoPrimary = Connect(mongoURI, "sample_analytics", PRIMARY, 1)
	defer mongoPrimary.Disconnect()
	defer mongoSecondary.Disconnect()
	assert.NotNil(t, mongoPrimary.Client)
	assert.NotNil(t, mongoSecondary.Client)

	sec := mongoSecondary.GetCollection("sample_analytics", "account")
	prim := mongoPrimary.GetCollection("sample_analytics", "account")

	assert.NotNil(t, sec)
	assert.NotNil(t, prim)
}
