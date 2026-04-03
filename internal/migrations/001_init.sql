-- 文档表
CREATE TABLE documents (
    id VARCHAR(64) PRIMARY KEY,
    content TEXT NOT NULL DEFAULT '',
    version BIGINT NOT NULL DEFAULT 0,
    created_by VARCHAR(64) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 操作记录表
CREATE TABLE operations (
    id VARCHAR(64) PRIMARY KEY,
    doc_id VARCHAR(64) NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    user_id VARCHAR(64) NOT NULL,
    type VARCHAR(16) NOT NULL,
    position INT NOT NULL,
    text TEXT,
    length INT,
    base_version BIGINT NOT NULL,
    version BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_doc_version (doc_id, version),
    INDEX idx_user (user_id),
    UNIQUE KEY uk_doc_version (doc_id, version)
);

-- 用户会话表
CREATE TABLE user_sessions (
    user_id VARCHAR(64) NOT NULL,
    doc_id VARCHAR(64) NOT NULL,
    role VARCHAR(16) NOT NULL DEFAULT 'viewer',
    last_seen_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, doc_id),
    INDEX idx_doc_online (doc_id, last_seen_at)
);

-- 文档快照表（用于持久化）
CREATE TABLE document_snapshots (
    doc_id VARCHAR(64) NOT NULL,
    content TEXT NOT NULL,
    version BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (doc_id, version),
    INDEX idx_version (version)
);