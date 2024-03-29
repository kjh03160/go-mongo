<p align="center">
  <a href="https://goreportcard.com/report/github.com/kjh03160/go-mongo"><img src="https://goreportcard.com/badge/github.com/kjh03160/go-mongo"></a>
</p>

# Golang Mongo Wrapper

The Wrapper that can be used in real projects based on [MongoDB Driver](https://github.com/mongodb/mongo-go-driver)

-------------------------
- [Goals](#the-goals-of-this-project)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
- [Feedback / Contribution](#feedback--contribution)
- [License](#license)

-------------------------
## The goals of this project
### Remove redundant codes
There are many duplicate codes when using mongo-driver.
For example, you may decode cursors on your custom struct and write every decoding cursor codes per struct.
This project will give you convenient code writing experience while using mongo-driver, by avoiding code writing like below.
- query timeout setting
  ```go
      ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
      defer cancel()
      cur, err := collection.Find(ctx, bson.D{})
      // ...
  ```
- decoding results 
  ```go
  cur, err := collection.Find(ctx, bson.D{})
  if err != nil { 
    // ...
  }
  defer cur.Close(ctx)
  
  for cur.Next(ctx) {
    var result CustomStruct
    err := cur.Decode(&result)
    if err != nil { 
        // ...
    }
  // ...
  }
  ```
### Log if slow query occurs
One of the important points when operating a server is to detect in advance if there is an abnormality, so that you could take action.
Usually the sign of database before failure is slow query. So, this project provide automatically slow query logging for you.
The Only thing you have to do is to make your logger implement our interface.

### Handle error easily
Mongo Driver provides some errors like duplicated key, not founded, and other internal errors.
However, as not found error is a variable, it should be compared with variable. 
On the other hand, duplicated key error is wrapped error, which means it should be unwrapped and need casting or catched by using provided function by Mongo Driver.
```go
  err = collection.FindOne(ctx, filter).Decode(&result)
  if err == mongo.ErrNoDocuments {
    fmt.Println("record does not exist")
  } else if mongo.IsDuplicateKeyError(err) {
    fmt.Println("record is duplicated")
  } else {
    fmt.Println("mongo internal error")
  }
 ``` 
As such, there is inconvenience in using it because each must be handled in a different way.
This inconvenience will also be resolved in this project.

-------------------------
## Requirements

- Go 1.19 or higher.
  - Wrapped functions use Generic of Golang.
- MongoDB version depends on MongoDB Driver version.
  - current version is `v1.11.1` which requires MongoDB 3.6 and higher.

-------------------------
## Installation

```bash
go get github.com/kjh03160/go-mongo
```

-------------------------
## Usage

### Connection
Pass Mongo Client Option instance you set, then you can get Mongo Client instance.
```go
import (
  "github.com/kjh03160/go-mongo/wrapper"
  "go.mongodb.org/mongo-driver/mongo/options"
)

func connect_mongo() {
  clientOptions := options.Client()
  clientOptions.ApplyURI(mongoURI).SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))
  
  mongoClient := wrapper.Connect(clientOptions)
}
```

### Create Collection
Call `NewCollection(client, databaseName, collectionName)` function with the type of struct to be decoded.
```go
import (
  "github.com/kjh03160/go-mongo/wrapper"
  "go.mongodb.org/mongo-driver/mongo/options"
)

type Account struct {
  AccountId int      `json:"account_id" bson:"account_id"`
  Limit     int      `json:"limit" bson:"limit"`
  Products  []string `json:"products" bson:"products"`
}

func getNewCollection() {
  mongoClient := wrapper.Connect(clientOptions)
  
  collection := wrapper.NewCollection[Account](mongoClient, "sample_analytics", "accounts")
}
```

If you have to decode bson result on different structs, you could do that by declaring each collection instance per struct.
For example,
```go
import (
  "github.com/kjh03160/go-mongo/wrapper"
  "go.mongodb.org/mongo-driver/mongo/options"
)

type Account struct {
  AccountId int      `json:"account_id" bson:"account_id"`
  Limit     int      `json:"limit" bson:"limit"`
  Products  []string `json:"products" bson:"products"`
}

type Product struct {
  ProductId int      `json:"product_id" bson:"product_id"`
  ProductName string `json:"product_name" bson:"product_name"`
}

func getNewCollection() {
  mongoClient := wrapper.Connect(clientOptions)
  
  account := mongo.NewCollection[Account](mongoClient, "sample_analytics", "accounts")
  product := mongo.NewCollection[Product](mongoClient, "sample_analytics", "accounts")
}
```

### Basic Query Usage
With collection instance, you can use wrapped query functions.
You should pass logger which implements our logger interface to log slow query.
```go
import (
  "github.com/kjh03160/go-mongo/wrapper"
)

func findOne() {
  collection := wrapper.NewCollection[Account](MongoClient, "sample_analytics", "accounts")
  
  // MyLogger is an example logger that implements logger interface.
  // The example will be described in Slow Query Usage.
  logger := MyLogger{logrus.New()}
  accountId := 1
  
  var t Account // declare variable with the type you set when you created collection
  // collection will find result, decode, and fill it in the variable you passed.
  // you should pass the pointer of variable.
  // if there is error, it will return error
  if err := collection.FindOne(&logger, &t, bson.M{"account_id": accountId}); err != nil {
    // The description of error types will be described in Error Handling.
    if errorType.IsNotFoundErr(err) {
      logger.Warn(err.Error())
      return
    }
      logger.Error(err.Error())
      return
  }
  logger.Info(t)
}

func findAll() {
  // findAll will return (result slice, err)
  // slice type is []T. T is the type of collection decoding struct you set
  all, err := collection.FindAll(&logger, bson.M{})
  if err != nil {
    if errorType.IsDecodeError(err) {
      logger.Warn("decode err", err.Error())
      return
    } else {
      logger.Error(err.Error())
      return
    }
  }
  logger.Info(all)
}

func insertMany()  {
  var result []Account
  accounts, _ := collection.FindAll(&logger, bson.M{})
  
  // example of insert slice
  var slice []interface{}
  for _, v := range accounts {
    slice = append(slice, v)
  }
  
  insertedIds, err := collection.InsertMany(&logger, slice)
  if err != nil {
    logger.Error(err.Error())
  }
  logger.Info(insertedIds)
}
```

### Transaction
Transaction function is provided for reducing redundant codes.
What you need to do is just pass session and transaction options, and transaction function.
```go
import (
  "github.com/kjh03160/go-mongo/wrapper"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

func transaction() {
  trx := func(sessCtx mongo.SessionContext) (interface{}, error) {
    // you can use transaction also by using WithTrx functions
    result, err := collection.FindAllWithTrx(&logger, bson.M{}, &sessCtx)
    
    // ...
  }

  trxOpt := options.TransactionOptions{}
  trxOpt.SetReadPreference(readpref.Primary())
  
  err := mongoClient.Transaction(&options.SessionOptions{}, &trxOpt, transaction)
  if err != nil {
	  // handle error
  } 
}
```

### Slow Query And Timeout
Also, if slow query is detected, logger will log about slow query info.
What you have to do is implementing `Logger` interface
```go
type Logger interface {
	// log slow query
	SlowQuery(msg string)
    // context cancel timeout setting -> ctx.Cancel() timeout
	// if mongo could not give the result until this time, 
	// wrapper cut the connection and return timeout error
	GetTimeoutDuration() time.Duration 
	// slow query time settings
	GetSlowQueryDurationOfOne() time.Duration
	GetSlowQueryDurationOfMany() time.Duration
	GetSlowQueryDurationOfBulk() time.Duration
	GetSlowQueryDurationOfAggregation() time.Duration
}
```
This is an example that implements `Logger` interface
```go

type MyLogger struct {
  *logrus.Logger
}

func (l *MyLogger) SlowQuery(msg string) {
  l.Error(msg)
}

func (l *MyLogger) GetTimeoutDuration() time.Duration {
  return 10 * time.Second
}

func (l *MyLogger) GetSlowQueryDurationOfOne() time.Duration {
  return 1 * time.Second
}

func (l *MyLogger) GetSlowQueryDurationOfMany() time.Duration {
  return 2 * time.Second
}

func (l *MyLogger) GetSlowQueryDurationOfBulk() time.Duration {
  return 3 * time.Second
}

func (l *MyLogger) GetSlowQueryDurationOfAggregation() time.Duration {
  return 10 * time.Second
}

func example() {
  // your logger
  logger := MyLogger{logrus.New()}
  // just pass it
  all, err := collection.FindAll(&logger, bson.M{})
}
```

### Error Handling
This project returns self-defined errors, not errors of Mongo Driver. And if error is `nil`, it guarantees database query is success
There are six errors we provide.
- `decodeError`
  - if `cursor.Decode()` provided by Mongo driver returns error.
- `notFoundError`
  - if `cursor.Decode()` provided by Mongo driver returns `mongo.ErrNoDocuments`
  - **In update query, if not matched count is zero in updateResult, it will return not found error.
    However, if ModifiedCount is zero, error is nil.**
  - In find all query, even if result slice is empty, it does not return `notFoundError`, but `nil`
- `duplicatedKeyError`
  - if Mongo Driver return error and `mongo.IsDuplicateKeyError(err)` is true, provided by Mongo Driver
- `timeoutError`
  - when context deadline exceed(`logger.GetTimeoutDuration()`) or `mongo.IsTimeout(err)` provided by Mongo Driver
- `mongoClientError`
  - an error during transaction session start
- `internalError`
  - all errors except the above

You can check which error occurs by using our functions.
```go
func IsDecodeError(err error) bool {}
func IsNotFoundErr(err error) bool {}
func IsDuplicatedKeyErr(err error) bool {}
func IsTimeoutError(err error) bool {}
func IsMongoClientError(err error) bool {}
// return true if error is one of internalError, timeoutError, mongoClientError
// Therefore, if you have to handle timeout or client error, you should filter them first with above function.
func IsDBInternalErr(err error) bool {}
```

If you filter error, then you could get error msg with `err.Error()`.
It provides you collection name, kind of error, and query info. (query info is provided only in query functions)
```text
accounts not found. | {query info: filter: map[account_id:0]}
```

-------------------------
## Feedback / Contribution

For help with the driver, please post on this repository issue.

-------------------------
## License

The Golang Mongo Wrapper is licensed under the [MIT License](LICENSE).
