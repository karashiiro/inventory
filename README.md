# inventory
Dockerized inventory microservice without authentication. Primarily intended to be used as a base for further development.

## Environment variables
```
INVENTORY_REDIS_LOCATION: Redis address
INVENTORY_PGL_CONNECTION_STRING: PostgreSQL connection string
INVENTORY_RMQ_CONNECTION_STRING: RabbitMQ connection string
INVENTORY_RMQ_CHANNEL: RabbitMQ channel name
```