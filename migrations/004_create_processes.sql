CREATE TABLE lcp_processes (
    id VARCHAR(36) PRIMARY KEY,
    status VARCHAR(32) NOT NULL,
    publication_id VARCHAR(36),
    error TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (publication_id) REFERENCES publications(id)
);
