# Firefly III Go Client Examples

This directory contains examples demonstrating how to use the Firefly III Go client library.

## Prerequisites

1. A running Firefly III instance
2. A personal access token from your Firefly III instance

## Running the Examples

To run these examples, you need to set the following environment variables:

```bash
export FIREFLY_BASE_URL="https://your-firefly-iii-instance/api"
export FIREFLY_ACCESS_TOKEN="your-personal-access-token"
```

Then you can run any example:

```bash
go run accounts.go
```

## Available Examples

### Account Operations (accounts.go)

This example demonstrates:
- Listing accounts
- Creating a new account
- Searching for accounts
- Updating account balances
- Retrieving account details
- Deleting accounts

### Transaction Operations (transactions.go)

This example demonstrates:
- Creating transactions
- Listing transactions
- Searching for transactions
- Retrieving transaction details
- Updating transactions
- Deleting transactions

## Notes

These examples are intended for educational purposes and should be modified for production use. In particular:

1. Error handling in these examples is minimal
2. No retry logic is implemented
3. The examples perform real operations against your Firefly III instance

Be cautious when running these examples against your production Firefly III instance, especially the ones that create, update, or delete data.

## Adjusting the Examples

You may need to adjust the examples to match your specific Firefly III setup:

1. Account types must be valid for your Firefly III instance
2. Currency codes should match currencies available in your instance
3. Category names may need to be changed to match existing categories or new ones you want to create 