# DOMAIN LAYER

- The place hold the core entities business can use anywhere in modules
- DON'T import any package in `adapter` like PostgreSQL, Redis, etc.
- DON'T import any package in `port` like Gin, etc.
