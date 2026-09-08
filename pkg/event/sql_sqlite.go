package event

import (
	"database/sql"
	"strings"

	watermillSQL "github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
)

// SQLiteSchema adapts Watermill SQL schema for SQLite database engine.
type SQLiteSchema struct {
	watermillSQL.DefaultMySQLSchema
}

// SchemaInitializingQueries returns table creation SQL for SQLite.
func (s SQLiteSchema) SchemaInitializingQueries(topic string) []watermillSQL.Query {
	createMessagesTable := strings.Join([]string{
		"CREATE TABLE IF NOT EXISTS " + s.MessagesTable(topic) + " (",
		"`offset` INTEGER PRIMARY KEY AUTOINCREMENT,",
		"`uuid` VARCHAR(36) NOT NULL,",
		"`created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,",
		"`payload` TEXT DEFAULT NULL,",
		"`metadata` TEXT DEFAULT NULL",
		");",
	}, "\n")

	return []watermillSQL.Query{{Query: createMessagesTable}}
}

// SubscribeIsolationLevel returns default isolation level for SQLite.
func (s SQLiteSchema) SubscribeIsolationLevel() sql.IsolationLevel {
	return sql.LevelDefault
}

// SQLiteOffsetsAdapter implements offset tracking for SQLite using ON CONFLICT clauses.
type SQLiteOffsetsAdapter struct {
	watermillSQL.DefaultMySQLOffsetsAdapter
}

// SchemaInitializingQueries returns offsets table creation SQL for SQLite.
func (a SQLiteOffsetsAdapter) SchemaInitializingQueries(topic string) []watermillSQL.Query {
	return []watermillSQL.Query{
		{
			Query: `
				CREATE TABLE IF NOT EXISTS ` + a.MessagesOffsetsTable(topic) + ` (
				consumer_group VARCHAR(255) NOT NULL,
				offset_acked BIGINT,
				offset_consumed BIGINT NOT NULL,
				PRIMARY KEY(consumer_group)
			)`,
		},
	}
}

// NextOffsetQuery returns query for current consumer group offset without unsupported locks.
func (a SQLiteOffsetsAdapter) NextOffsetQuery(topic, consumerGroup string) watermillSQL.Query {
	return watermillSQL.Query{
		Query: `SELECT COALESCE(
				(SELECT offset_acked
				 FROM ` + a.MessagesOffsetsTable(topic) + `
				 WHERE consumer_group=?
				), 0)`,
		Args: []any{consumerGroup},
	}
}

// AckMessageQuery updates acked offset on SQLite using ON CONFLICT.
func (a SQLiteOffsetsAdapter) AckMessageQuery(topic string, row watermillSQL.Row, consumerGroup string) watermillSQL.Query {
	ackQuery := `INSERT INTO ` + a.MessagesOffsetsTable(topic) + ` (offset_consumed, offset_acked, consumer_group)
		VALUES (?, ?, ?)
		ON CONFLICT(consumer_group) DO UPDATE SET offset_consumed=excluded.offset_consumed, offset_acked=excluded.offset_acked`

	return watermillSQL.Query{Query: ackQuery, Args: []any{row.Offset, row.Offset, consumerGroup}}
}

// ConsumedMessageQuery updates consumed offset on SQLite using ON CONFLICT.
func (a SQLiteOffsetsAdapter) ConsumedMessageQuery(topic string, row watermillSQL.Row, consumerGroup string, consumerULID []byte) watermillSQL.Query {
	consumedQuery := `INSERT INTO ` + a.MessagesOffsetsTable(topic) + ` (offset_consumed, consumer_group)
		VALUES (?, ?)
		ON CONFLICT(consumer_group) DO UPDATE SET offset_consumed=excluded.offset_consumed`
	return watermillSQL.Query{Query: consumedQuery, Args: []any{row.Offset, consumerGroup}}
}
