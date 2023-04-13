# APPLICATION LAYER

- The place hold the logic business handling
- Using CQRS patterns
  - Query: hold features get data from the database
  - Command: hold features modify data to the database
- DON'T import any package in `adapter` like PostgreSQL, Redis, etc.
- DON'T import any package in `port` like Gin, etc.
