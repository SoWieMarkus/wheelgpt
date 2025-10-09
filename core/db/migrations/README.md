# Database Migrations

This directory contains database migration files for the WheelGPT application. Migrations are used to manage changes to the database schema over time, ensuring that the database structure remains consistent with the application's requirements.

Migrations should follow this naming convention:

```txt
YYYY-MM-DD_description.sql
```

`YYYY-MM-DD` is the date of the migration in the format `YYYY-MM-DD`, and `description` is a brief description of the migration's purpose. Example `2025-01-15_add_channel_settings.sql`.

Each migration file should contain valid SQL statements that define the changes to be applied to the database schema. When adding a new migration, ensure that it is idempotent and can be applied safely without causing errors if run multiple times.

Do not create tables with sql statements. All new tables are automatically created from the models in `core/db/models`. Only use migrations for altering existing tables, adding indexes, or other schema modifications that cannot be handled by the ORM.
