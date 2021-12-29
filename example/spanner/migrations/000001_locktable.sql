CREATE TABLE ResourceLockTable (
    Name STRING(MAX) NOT NULL,
    Heartbeat TIMESTAMP OPTIONS (allow_commit_timestamp=true),
    Token TIMESTAMP OPTIONS (allow_commit_timestamp=true),
    Writer STRING(MAX),
) PRIMARY KEY (Name)