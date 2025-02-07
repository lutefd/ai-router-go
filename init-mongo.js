db.createCollection("users");
db.createCollection("chats");

db.users.createIndex({ email: 1 }, { unique: true });
db.users.createIndex({ id: 1 }, { unique: true });
db.chats.createIndex({ id: 1 }, { unique: true });
