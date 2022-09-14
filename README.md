# Learn

This project contains all code for the Johnny Morrice programmer teaching course.

## Learnblog

[Learnblog](/learnblog) is a web application for teaching golang and development concepts.

The application is a simple blog engine.  In the course we will add various features, and deploy it to the cloud.

The stack we will implement is:

* Go
* Docker
* Postgres
* Redis
* Solidjs
* Render.com
* Netlify

## Useful commands

### Run a migration:

```
docker run --network learnblog_learnblog-network --env-file env/dev/learnblog.env -it --rm learnblog /app/entrypoint.sh --command migrate-up
```

### Open the database

```
docker exec -it $DB_CONTAINER_ID /usr/local/bin/psql -U postgres -d postgres
```

