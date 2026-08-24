# Earthbound

## Postgres User Setup

### Local

For local testing i set the user up with the following, do not follow this for a production setup:

```sql
CREATE DATABASE earthbound;
CREATE USER earthbound WITH PASSWORD 'earthbound';
GRANT ALL PRIVILEGES ON DATABASE earthbound TO earthbound;
```
