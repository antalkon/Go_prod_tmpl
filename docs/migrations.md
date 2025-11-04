# Миграции

Файлы миграций

- Хранятся в каталоге `./migrations`.
- Каждая миграция состоит минимум из двух файлов: `NNNN_description.up.sql` и `NNNN_description.down.sql`.
- `NNNN` — числовой префикс (например `0001`, `0002`) для упорядочивания.

Пример:

```
0001_create_pings_table.up.sql
0001_create_pings_table.down.sql
```

Создание новой миграции

1. Создайте две SQL-ручки (`.up.sql` и `.down.sql`) с последовательным префиксом и понятным именем.
2. В `.up.sql` — DDL/DML для применения, в `.down.sql` — откат

Пример содержимого `0001_create_pings_table.up.sql`:

```sql
CREATE TABLE IF NOT EXISTS pings (
	id BIGSERIAL PRIMARY KEY,
	message TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

И соответствующий `0001_create_pings_table.down.sql`:

```sql
DROP TABLE IF EXISTS pings;
```

Запуск миграций (local)

cli инструмент — `golang-migrate`  для контролируемого применения миграций.

cli tmpl:

```bash
# применить все доступные миграции
migrate -path ./migrations -database "$DATABASE_DSN" up

# откатить одну миграцию
migrate -path ./migrations -database "$DATABASE_DSN" down 1

# привести БД к конкретной версии
migrate -path ./migrations -database "$DATABASE_DSN" goto 3
```

`$DATABASE_DSN` tml:

```
postgres://user:pass@localhost:5432/dbname?sslmode=disable
```

Автомиграции в приложении

Для удобства  сделал пакет `pkg/migrations` и флаг `AUTO_MIGRATE`.
При включённом `AUTO_MIGRATE=true` приложение вызывает `pkg/migrations.MigrateUp("./migrations", cfg.DatabaseDSN, log)` при старте.


пример:

```go
if cfg.AutoMigrate {
		if err := migrations.MigrateUp("./migrations", cfg.DatabaseDSN, log); err != nil {
				return err
		}
}
```

Откат миграций программно

В пакете `pkg/migrations` есть функция `MigrateDown(migrationsDir, databaseURL, log)` — она откатывает шаг назад.

