# DynamoDB データ移行ガイド (シングルテーブル設計)

## テーブル

テーブル名: `MathOvercome`

| キー | 型 | 説明 |
|---|---|---|
| `pk` | String | パーティションキー |
| `sk` | String | ソートキー |
| `gsi1pk` | String | GSI1 パーティションキー |
| `gsi1sk` | String | GSI1 ソートキー |

### GSI

| インデックス名 | キー | 用途 |
|---|---|---|
| GSI1 | gsi1pk (HASH), gsi1sk (RANGE) | カテゴリ別問題一覧、ユーザーのセッション一覧など |
| GSI2 | user_id (HASH) | ログイン文字列でのユーザー検索 (FindUserByUserID) |

## アイテムキーパターン

| エンティティ | pk | sk | gsi1pk | gsi1sk |
|---|---|---|---|---|
| CATEGORY | `CATEGORY#<id>` | `#METADATA` | `CATEGORY` | `CATEGORY#<id>` |
| PROBLEM | `PROBLEM#<id>` | `#METADATA` | `CATEGORY#<cat_id>` | `PROBLEM#<id>` |
| CHOICE | `PROBLEM#<problem_id>` | `CHOICE#<id>` | (なし) | (なし) |
| USER | `USER#<id>` | `#METADATA` | `USER` | `USER#<id>` |
| TESTSESSION | `SESSION#<id>` | `#METADATA` | `USER#<user_id>` | `SESSION#<id>` |
| SESSIONPROBLEM | `SESSION#<session_id>` | `SP#<id>` | (なし) | (なし) |

## アクセスパターン

| 操作 | クエリ |
|---|---|
| カテゴリ一覧 | GSI1: gsi1pk = `CATEGORY` |
| カテゴリ別問題一覧 | GSI1: gsi1pk = `CATEGORY#1` |
| 問題＋選択肢取得 | PK: `PROBLEM#197`, sk begins_with `CHOICE#` (or `#METADATA`) |
| ユーザー一覧 | GSI1: gsi1pk = `USER` |
| ユーザーのセッション一覧 | GSI1: gsi1pk = `USER#2` |
| セッションの解答一覧 | PK: `SESSION#143`, sk begins_with `SP#` |

## テーブル作成

```bash
bash create_tables.sh
```

## データアップロード

テーブル作成後、以下を実行してください：

```bash
bash upload.sh
```

> **注意**: `create_tables.sh` と `upload.sh` 内の `REGION` 変数を
> 使用するAWSリージョンに合わせて変更してください（デフォルト: ap-northeast-1）。

## データファイル構成

```
data/
  data_01.json   (25件: CATEGORY×7 + PROBLEM 197-214)
  data_02.json   (25件: PROBLEM 215-238 + CHOICE 487)
  data_03.json   (25件: CHOICE 488-512)
  data_04.json   (25件: CHOICE 513-537)
  data_05.json   (25件: CHOICE 538-562)
  data_06.json   (25件: CHOICE 563-587)
  data_07.json   (25件: CHOICE 588-612)
  data_08.json   (25件: CHOICE 613-637)
  data_09.json   (25件: CHOICE 638-654 + USER×5 + SESSION 143-145)
  data_10.json   (25件: SESSION 146-153 + SESSIONPROBLEM 1715-1731)
  data_11.json   (25件: SESSIONPROBLEM 1732-1756)
  data_12.json   (25件: SESSIONPROBLEM 1757-1781)
  data_13.json   (25件: SESSIONPROBLEM 1782-1806)
  data_14.json   (25件: SESSIONPROBLEM 1807-1831)
  data_15.json   (17件: SESSIONPROBLEM 1832-1848)
```

合計: 367件

## アイテム属性一覧

### CATEGORY
| 属性 | 型 | 備考 |
|---|---|---|
| pk | String | `CATEGORY#<id>` |
| sk | String | `#METADATA` |
| gsi1pk | String | `CATEGORY` |
| gsi1sk | String | `CATEGORY#<id>` |
| id | Number | |
| name | String | |

### PROBLEM
| 属性 | 型 | 備考 |
|---|---|---|
| pk | String | `PROBLEM#<id>` |
| sk | String | `#METADATA` |
| gsi1pk | String | `CATEGORY#<category_id>` |
| gsi1sk | String | `PROBLEM#<id>` |
| id | Number | |
| category_id | Number | |
| question | String | HTML含む可 |
| hint | String | |

### CHOICE
| 属性 | 型 | 備考 |
|---|---|---|
| pk | String | `PROBLEM#<problem_id>` |
| sk | String | `CHOICE#<id>` |
| id | Number | |
| problem_id | Number | |
| choice_text | String | |
| is_correct | Boolean | |

### USER
| 属性 | 型 | 備考 |
|---|---|---|
| pk | String | `USER#<id>` |
| sk | String | `#METADATA` |
| gsi1pk | String | `USER` |
| gsi1sk | String | `USER#<id>` |
| id | Number | |
| user_name | String | |
| user_id | String | ユニーク |
| password | String | ⚠️ 平文パスワード — 本番環境ではハッシュ化を推奨 |

### TESTSESSION
| 属性 | 型 | 備考 |
|---|---|---|
| pk | String | `SESSION#<id>` |
| sk | String | `#METADATA` |
| gsi1pk | String | `USER#<user_id>` |
| gsi1sk | String | `SESSION#<id>` |
| id | Number | |
| user_id | Number | |
| include_integers | Boolean | 整数問題を含むか |
| start_time | String | datetime文字列 |

### SESSIONPROBLEM
| 属性 | 型 | 備考 |
|---|---|---|
| pk | String | `SESSION#<session_id>` |
| sk | String | `SP#<id>` |
| id | Number | |
| session_id | Number | |
| problem_id | Number | |
| selected_choice_id | Number | 未回答時は属性なし |
| is_correct | Boolean | 未回答時は属性なし |
