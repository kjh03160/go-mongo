package wrapper

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Test_Connection(t *testing.T) {
	_ = os.Getenv("user")
	_ = os.Getenv("password")

	t.Run("connection success", func(t *testing.T) {
		var mongoSecondary *Client
		var mongoPrimary *Client

		mongoSecondary = Connect(getMongoConfig())
		mongoPrimary = Connect(getMongoConfig())
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
		clientOptions := options.Client()
		clientOptions.ApplyURI("mongoURI").SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))
		mongoSecondary = Connect(clientOptions)
		assert.Nil(t, mongoSecondary.Client)
	})

	t.Run("connection failed - password", func(t *testing.T) {
		var mongoSecondary *Client
		defer func() {
			r := recover()
			assert.NotNil(t, r)
		}()
		mongoSecondary = Connect(&options.ClientOptions{})
		assert.Nil(t, mongoSecondary.Client)
	})
}

func Test_GetDatabase(t *testing.T) {
	var mongoSecondary *Client
	var mongoPrimary *Client
	clientOptions := getMongoConfig()
	clientOptions.SetReadPreference(readpref.SecondaryPreferred())
	mongoSecondary = Connect(clientOptions)
	clientOptions2 := getMongoConfig()
	clientOptions2.SetReadPreference(readpref.PrimaryPreferred())
	mongoPrimary = Connect(clientOptions2)
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
	var mongoSecondary *Client
	var mongoPrimary *Client
	clientOptions := getMongoConfig()
	clientOptions.SetReadPreference(readpref.SecondaryPreferred())
	mongoSecondary = Connect(clientOptions)
	clientOptions2 := getMongoConfig()
	clientOptions2.SetReadPreference(readpref.PrimaryPreferred())
	mongoPrimary = Connect(clientOptions2)
	defer mongoPrimary.Disconnect()
	defer mongoSecondary.Disconnect()
	assert.NotNil(t, mongoPrimary.Client)
	assert.NotNil(t, mongoSecondary.Client)

	sec := mongoSecondary.GetCollection("sample_analytics", "account")
	prim := mongoPrimary.GetCollection("sample_analytics", "account")

	assert.NotNil(t, sec)
	assert.NotNil(t, prim)
}
