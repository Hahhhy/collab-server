//文档表
CREATE TABLE documents (
    //SERIAL类型会自动递增，适合用作主键
    id SERIAL PRIMARY KEY,
    //VARCHAR(255)表示字符串类型，NOT NULL表示不能为空
    title VARCHAR(255) NOT NULL,
    //TEXT类型适合存储较长的文本内容，NOT NULL表示不能为空
    content TEXT NOT NULL,
    //TIMESTAMP类型适合存储时间戳，DEFAULT CURRENT_TIMESTAMP表示默认值为当前时间
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

--操作记录表
CREATE TABLE operation_logs (
    id SERIAL PRIMARY KEY,
    doc_id INTEGER NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    operation VARCHAR(255) NOT NULL,
    operator VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (doc_id) REFERENCES documents(id) ON DELETE CASCADE
);

--用户会话表
CREATE TABLE user_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    session_token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL
);

--文档快照表
CREATE TABLE document_snapshots (
    id SERIAL PRIMARY KEY,
    doc_id INTEGER NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    snapshot TEXT NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);