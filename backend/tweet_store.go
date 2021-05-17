package backend

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

// NewTweetStore is New TweetStore
func NewTweetStore(sc *spanner.Client) *TweetStore {
	return &TweetStore{
		sc: sc,
	}
}

// Tweet is TweetTable Row
type Tweet struct {
	ID             string `spanner:"Id"`
	Author         string
	Content        string
	Count          int64
	Favos          []string
	Sort           int64
	ShardCreatedAt int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	CommitedAt     time.Time
	SchemaVersion  int64
}

type TweetStore struct {
	sc *spanner.Client
}

func (s *TweetStore) TableName() string {
	return "Tweet"
}

// InsertChain is Mutation APIを利用して、複数RowのInsertを実行
func (s *TweetStore) InsertChain(ctx context.Context, content string) ([]string, error) {
	ctx, span := startSpan(ctx, "insertChain")
	defer span.End()

	now := time.Now()
	var ms []*spanner.Mutation
	var ids []string

	for i := 0; i < 3; i++ {
		id := uuid.New().String()
		ids = append(ids, id)
		m, err := spanner.InsertStruct(s.TableName(), &Tweet{
			ID:         uuid.New().String(),
			Author:     "srun",
			Content:    fmt.Sprintf("%s-%d", content, i),
			Favos:      []string{},
			CreatedAt:  now,
			UpdatedAt:  now,
			CommitedAt: spanner.CommitTimestamp,
		})
		if err != nil {
			return ids, err
		}
		ms = append(ms, m)
	}

	_, err := s.sc.ReadWriteTransaction(ctx, func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
		return tx.BufferWrite(ms)
	})
	if err != nil {
		return ids, err
	}
	return ids, nil
}

// UpdateChain is Mutation APIを利用して、複数RowのUpdateを実行
func (s *TweetStore) UpdateChain(ctx context.Context, ids []string, content string) error {
	ctx, span := startSpan(ctx, "updateChain")
	defer span.End()

	now := time.Now()
	var ms []*spanner.Mutation
	for _, id := range ids {
		ids = append(ids, id)
		m := spanner.UpdateMap(s.TableName(), map[string]interface{}{
			"ID":         id,
			"Author":     "srun",
			"Content":    content,
			"UpdatedAt":  now,
			"CommitedAt": spanner.CommitTimestamp,
		})
		ms = append(ms, m)
	}

	_, err := s.sc.ReadWriteTransaction(ctx, func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
		return tx.BufferWrite(ms)
	})
	if err != nil {
		return err
	}
	return nil
}

// UpdateDMLChain is DMLを利用して、複数RowのUpdateを実行
func (s *TweetStore) UpdateDMLChain(ctx context.Context, ids []string, content string) ([]string, error) {
	ctx, span := startSpan(ctx, "updateDMLChain")
	defer span.End()

	now := time.Now()
	_, err := s.sc.ReadWriteTransaction(ctx, func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
		for _, id := range ids {
			stmt := spanner.Statement{
				SQL: `UPDATE Tweet SET Content = @Content, UpdatedAt = @UpdatedAt WHERE Id = @Id`,
			}
			stmt.Params = map[string]interface{}{
				"UpdatedAt": now,
				"Id":        id,
				"Content":   content,
			}
			_, err := tx.Update(ctx, stmt)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return ids, err
	}
	return ids, nil
}

// UpdateAndSelect is Mutation APIを利用して、Updateを間に強引に挟んでみると、反映されるのかな？
func (s *TweetStore) UpdateAndSelect(ctx context.Context, id string) (int64, error) {
	ctx, span := startSpan(ctx, "updateAndSelectChain")
	defer span.End()

	now := time.Now()

	var count int64
	_, err := s.sc.ReadWriteTransaction(ctx, func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
		row, err := tx.ReadRow(ctx, s.TableName(), spanner.Key{id}, []string{"Count"})
		if err != nil {
			return err
		}
		if err := row.ColumnByName("Count", &count); err != nil {
			return err
		}
		fmt.Printf("Before Count %d\n", count)

		var ms []*spanner.Mutation
		mu := spanner.UpdateMap(s.TableName(), map[string]interface{}{
			"ID":         id,
			"Count":      count + 1,
			"UpdatedAt":  now,
			"CommitedAt": spanner.CommitTimestamp,
		})
		ms = append(ms, mu)
		if err := tx.BufferWrite(ms); err != nil {
			return err
		}

		row, err = tx.ReadRow(ctx, s.TableName(), spanner.Key{id}, []string{"Count"})
		if err != nil {
			return err
		}
		if err := row.ColumnByName("Count", &count); err != nil {
			return err
		}
		fmt.Printf("After Count %d\n", count)

		return nil
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// UpdateDMLAndSelect is DMLを利用してUpdateを行い、同じTxでReadをする
//
// DMLはtx.Update()を呼んだ時点で実行されるので、同じTxの以降のReadで更新が反映される
func (s *TweetStore) UpdateDMLAndSelect(ctx context.Context, id string) (int64, error) {
	ctx, span := startSpan(ctx, "updateDMLAndSelectChain")
	defer span.End()

	var count int64
	now := time.Now()
	_, err := s.sc.ReadWriteTransaction(ctx, func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
		{
			iter := tx.Query(ctx, spanner.Statement{
				SQL: `SELECT Count FROM Tweet WHERE Id = @Id`,
				Params: map[string]interface{}{
					"Id": id,
				},
			})
			defer iter.Stop()
			for {
				row, err := iter.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					return err
				}
				if err := row.ColumnByName("Count", &count); err != nil {
					return err
				}
			}
			fmt.Printf("Before Count %d\n", count)
		}
		{
			stmt := spanner.Statement{
				SQL: `UPDATE Tweet SET Count = Count + 1, UpdatedAt = @UpdatedAt WHERE Id = @Id`,
				Params: map[string]interface{}{
					"UpdatedAt": now,
					"Id":        id,
				},
			}
			// ここでUPDATE文が実行され、反映されるので、以降の処理でCountを再びSpannerから取得した場合は +1 されている
			_, err := tx.Update(ctx, stmt)
			if err != nil {
				return err
			}
		}
		{
			iter := tx.Query(ctx, spanner.Statement{
				SQL: `SELECT Count FROM Tweet WHERE Id = @Id`,
				Params: map[string]interface{}{
					"Id": id,
				},
			})
			defer iter.Stop()
			for {
				row, err := iter.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					return err
				}
				if err := row.ColumnByName("Count", &count); err != nil {
					return err
				}
			}
			fmt.Printf("After Count %d\n", count)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// UpdateBatchDMLChain is 複数のDMLを一度に実行するtx.BatchUpdate()のサンプル
func (s *TweetStore) UpdateBatchDMLChain(ctx context.Context, ids []string, content string) ([]string, error) {
	ctx, span := startSpan(ctx, "updateBatchDMLChain")
	defer span.End()

	var stmts []spanner.Statement
	now := time.Now()
	for _, id := range ids {
		stmt := spanner.Statement{
			SQL: `UPDATE Tweet SET Content = @Content, UpdatedAt = @UpdatedAt WHERE Id = @Id`,
		}
		stmt.Params = map[string]interface{}{
			"UpdatedAt": now,
			"Id":        id,
			"Content":   content,
		}
		stmts = append(stmts, stmt)
	}
	_, err := s.sc.ReadWriteTransaction(ctx, func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
		_, err := tx.BatchUpdate(ctx, stmts)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return ids, err
	}
	return ids, nil
}
