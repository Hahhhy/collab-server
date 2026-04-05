package storage

import (
	"TO/internal/models"
	"database/sql"
	"time"
)

type Storage struct {
	db *sql.DB
}

// NewStorage 创建一个新的 Storage 实例
func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

// storage 层
func (s *Storage) GetUserRole(userID, docID string) (models.UserRole, error) {
	query := `SELECT role FROM user_sessions WHERE user_id = $1 AND doc_id = $2`
	var role models.UserRole
	err := s.db.QueryRow(query, userID, docID).Scan(&role)
	return role, err
}

func (s *Storage) GetUserConnectionStatus(userID string) (models.ConnectionStatus, error) {
	query := `SELECT state FROM user_sessions WHERE user_id = $1`
	var state models.ConnectionStatus
	err := s.db.QueryRow(query, userID).Scan(&state)
	return state, err
}

// SaveSnapshot 保存文档快照
func (s *Storage) SaveSnapshot(docID string, content string, version int) error {
	query := `
        INSERT INTO document_snapshots (doc_id, content, version, created_at)
        VALUES ($1, $2, $3, NOW())
    `
	_, err := s.db.Exec(query, docID, content, version)
	return err
}

// GetLatestSnapshot 获取文档的最新快照
func (s *Storage) GetLatestSnapshot(docID string) (*models.DocumentSnapshot, error) {
	query := `
        SELECT doc_id, content, version, created_at
        FROM document_snapshots
        WHERE doc_id = $1
        ORDER BY version DESC
        LIMIT 1
    `
	var snap models.DocumentSnapshot
	var createdAt time.Time
	err := s.db.QueryRow(query, docID).Scan(&snap.DocID, &snap.Content, &snap.Version, &createdAt)
	if err == sql.ErrNoRows {
		// 没有快照，返回空快照（文档为空）
		return &models.DocumentSnapshot{
			DocID:   docID,
			Content: "",
			Version: 0,
		}, nil
	}
	if err != nil {
		return nil, err
	}
	return &snap, nil
}

// GetOperationsSince 获取文档在某个版本之后的所有操作（按版本升序）
func (s *Storage) GetOperationsSince(docID string, sinceVersion int) ([]models.Operation, error) {
	query := `
        SELECT id, doc_id, user_id, type, position, text, length, base_version, version, created_at
        FROM operations
        WHERE doc_id = $1 AND version > $2
        ORDER BY version ASC
    `
	rows, err := s.db.Query(query, docID, sinceVersion)
	if err != nil {
		return nil, err
	}
	//这里一定要记得关闭rows，否则会有内存泄漏，导致连接池耗尽，无法继续查询数据库
	//Query 用于执行可能返回多行的查询，它会返回一个 *sql.Rows 结果集，
	// 需要通过 rows.Next() 来迭代结果集中的每一行，
	// 并且在使用完毕后调用 rows.Close() 来释放资源。
	//底层独占一个数据库连接用于拉取结果集
	defer rows.Close()

	var ops []models.Operation
	for rows.Next() {
		var op models.Operation
		err := rows.Scan(
			&op.ID, &op.DocID, &op.UserID, &op.Type,
			&op.Position, &op.Text, &op.Length,
			&op.BaseVersion, &op.Version, &op.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		ops = append(ops, op)
	}
	return ops, nil
}

// SaveOperation 保存一个操作
func (s *Storage) SaveOperation(op models.Operation) error {
	query := `
        INSERT INTO operations (id, doc_id, user_id, type, position, text, length, base_version, version, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
    `
	_, err := s.db.Exec(query, op.ID, op.DocID, op.UserID, op.Type, op.Position, op.Text, op.Length, op.BaseVersion, op.Version)
	return err
}
