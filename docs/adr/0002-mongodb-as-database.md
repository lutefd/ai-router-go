# 2. MongoDB as Primary Database

## Status

Accepted

## Context

We needed a database solution that could handle:

- Flexible schema for chat messages with varying content types
- High-performance read/write operations for real-time chat
- Easy scalability for growing data
- Support for real-time operations
- Efficient querying of chat history
- Simple data modeling for message threads

We evaluated:

- PostgreSQL
- MySQL
- MongoDB
- Redis
- CockroachDB

## Decision

We chose MongoDB as our primary database for the following reasons:

1. **Data Model Fit**

   ```json
   {
   	"_id": "chat_123",
   	"user": "user_456",
   	"title": "Chat Title",
   	"messages": [
   		{
   			"_id": "msg_789",
   			"text": "Hello",
   			"role": "user",
   			"ai": "gemini",
   			"sent_at": "2024-03-10T12:00:00Z"
   		}
   	],
   	"created_at": "2024-03-10T12:00:00Z",
   	"updated_at": "2024-03-10T12:00:00Z"
   }
   ```

2. **Query Capabilities**

   ```go
   filter := bson.M{
       "user": userID,
       "created_at": bson.M{
           "$gte": startDate,
           "$lte": endDate,
       },
   }
   ```

3. **Indexing Strategy**
   ```go
   indexes := []mongo.IndexModel{
       {
           Keys: bson.D{{Key: "user", Value: 1}, {Key: "created_at", Value: -1}},
           Options: options.Index().SetUnique(false),
       },
       {
           Keys: bson.D{{Key: "_id", Value: 1}},
           Options: options.Index().SetUnique(true),
       },
   }
   ```

## Consequences

### Positive

- Schema flexibility for varying message formats
- Good performance for read/write operations
- Native support for JSON-like documents
- Easy horizontal scaling
- Rich query capabilities
- Built-in support for geospatial queries
- Good driver support for Go
- Simple backup and restore

### Negative

- No built-in support for ACID transactions
- Higher memory usage compared to SQL databases
- Requires careful index planning
- Limited join capabilities
- No built-in schema validation

### Mitigations

1. For data consistency:

   ```go
   session, err := client.StartSession()
   if err != nil {
       return err
   }
   defer session.EndSession(context.Background())

   callback := func(sessionContext mongo.SessionContext) error {
       return nil
   }

   err = session.WithTransaction(context.Background(), callback)
   ```

2. For schema validation:

   ```go
   validator := bson.M{
       "$jsonSchema": bson.M{
           "required": []string{"user", "title", "messages"},
           "properties": bson.M{
               "user": bson.M{"type": "string"},
               "title": bson.M{"type": "string"},
               "messages": bson.M{"type": "array"},
           },
       },
   }
   ```

3. For indexing:
   ```go
   opts := options.CreateCollection().SetValidator(validator)
   err = db.CreateCollection(context.Background(), "chats", opts)
   ```
