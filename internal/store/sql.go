package store

const (
	postgresWrite = `
  INSERT INTO notifications
    (recipients, message, tags, priority, type, status, created_at, updated_at)
  VALUES
    (:recipients, :message, :tags, :priority, :type, :status, :created_at, :updated_at);
  `
	postgresRetrieveByIDSQL = `
    SELECT id, recipients, message, tags, priority, type, status, created_at, updated_at
    FROM notifications
    WHERE id = $1;
  `
)
