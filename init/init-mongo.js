var db = db.getSiblingDB('my_mongodb');

// Create the collection for BlockTransactions
db.createCollection('BlockTransactions', {
    validator: {
        $jsonSchema: {
            bsonType: 'object',
            required: ['block_num', 'transactions'],
            properties: {
                block_num: { bsonType: 'int' },
                transactions: {
                    bsonType: 'array',
                    items: { bsonType: 'string' }
                }
            }
        }
    },
    validationAction: 'error',
    validationLevel: 'strict'
});

db.BlockTransactions.createIndex({ block_num: 1 }, { unique: true });

// Create the collection for Transation
db.createCollection('Transaction', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['tx_hash', 'from', 'to', 'nonce', 'data', 'value', 'logs'],
      properties: {
        tx_hash: {
          bsonType: 'string',
          description: 'must be a string and is required',
        },
        from: {
          bsonType: 'string',
          description: 'must be a string and is required',
        },
        to: {
          bsonType: 'string',
          description: 'must be a string and is required',
        },
        nonce: {
          bsonType: 'int',
          description: 'must be an integer and is required',
        },
        data: {
          bsonType: 'string',
          description: 'must be a string and is required',
        },
        value: {
          bsonType: 'string',
          description: 'must be a string and is required',
        },
        logs: {
          bsonType: 'array',
          description: 'must be an array and is required',
          items: {
            bsonType: 'object',
            required: ['index', 'data'],
            properties: {
              index: {
                bsonType: 'int',
                description: 'must be an integer and is required',
              },
              data: {
                bsonType: 'string',
                description: 'must be a string and is required',
              },
            },
          },
        },
      },
    },
  },
  validationLevel: 'strict',
  validationAction: 'error',
});

db.Transaction.createIndex({ tx_hash: 1 });
