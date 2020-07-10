package sqlmock

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSQLMock_Select(t *testing.T) {

	// モックDBの初期化
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to init db mock")
	}
	defer db.Close()

	id := 1

	// dbドライバに対する操作のモック定義
	columns := []string{"id", "title", "content"}
	mock.ExpectQuery("SELECT (.+) FROM ARTICLES"). // expectedSQL: 想定される実行クエリをregexpで指定（指定文字列が含まれるかどうかを見る）
		WithArgs(id).                                          // 想定されるプリペアドステートメントへの引数
		WillReturnRows(sqlmock.NewRows(columns).AddRow(1, "test title", "test content")) // 返戻する行情報の指定

	// テスト対象関数call
	article, err := GetByID(id, db)
	if err != nil {
		t.Fatalf("failed to get article: %s", err)
	}
	fmt.Printf("%v", article)

	// mock定義の期待操作が順序道理に実行されたか検査
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("failed to ExpectationWerMet(): %s", err)
	}
}

func TestSQLMock_Insert(t *testing.T) {

	// モックDBの初期化
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to init db mock")
	}
	defer db.Close()

	id := 1
	title := "test title"
	content := "test content"

	// dbドライバに対する操作のモック定義
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO ARTICLES"). // 想定される実行SQLをregexpで指定（指定文字列が含まれるかどうかを見る）
		WithArgs(id, title, content).                    // 想定されるプリペアドステートメントへの引数
		WillReturnResult(sqlmock.NewResult(1, 1)) // 想定されるExec関数の結果を指定
	mock.ExpectCommit()

	// テスト対象関数call
	if err := Create(id, title, content, db); err != nil {
		t.Fatalf("failed to get article: %s", err)
	}

	// mock定義の期待操作が順序道理に実行されたか検査
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("failed to ExpectationWerMet(): %s", err)
	}
}
