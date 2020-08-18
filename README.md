# vtgqlgen
Variant of vtgql but with dataloaders gqlgen tool

# Problem:
1. How to solve the N+1 queries problem?
1. Cleaning up the object models so that they represent their attributes in the respective datasources.
1. Using multiple datasources and of different kinds (ex: DB, REST APIs)

# Reference
github.com/varunturlapati/vtgql

# Tools
1. `gqlgen`
1. `dataloaden`
1. `sqlc` - currently doesn't support much more than Postgres. (This project doesn't use `sqlc`
