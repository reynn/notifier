CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE notification_type as ENUM ('HIGH','LOW','DEFAULT');
CREATE TYPE notification_priority as ENUM ('EMAIL','SMS','WEBHOOK','PUSH','UNSET');
CREATE TYPE notification_status as ENUM ('SUBMITTED','COMPLETED','FAILED');

CREATE TABLE notifications (
    id         uuid                  primary key default uuid_generate_v4(),
  	recipients text[]                not     null,
    message    bytea                 not     null,
    tags       text[],
    priority   notification_priority not     null,
    type       notification_type     not     null,
    status     notification_status   not     null,
    created_at timestamptz           DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz
);
