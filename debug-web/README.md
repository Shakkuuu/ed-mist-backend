# Mist ED Backend Debug Web Tool

React + TypeScript + Tailwind CSS を使用したWeb版デバッグツールです。

## 機能

- データの一覧表示（組織、ユーザー、部屋、デバイス、科目、授業）
- テストデータの作成・削除
- シードデータの一括作成
- データベースのリセット
- リアルタイムでのデータ更新

## セットアップ

### 1. 依存関係のインストール

```bash
npm install
```

### 2. 環境変数の設定

```bash
# .env.local ファイルを作成
echo "REACT_APP_API_URL=http://localhost:8080" > .env.local
```

### 3. 開発サーバーの起動

```bash
npm start
```

ブラウザで `http://localhost:3000` にアクセスしてください。

## 使用方法

### 基本的な操作フロー

1. **データベースリセット**: 全データを削除してクリーンな状態にする
2. **シードデータ作成**: 開発・テスト用のサンプルデータを一括作成
3. **データ確認**: 各テーブルでデータが正しく作成されていることを確認
4. **個別操作**: 必要に応じて個別のデータ作成・削除を実行

### 画面構成

- **コントロールパネル**: シードデータ作成・データベースリセット
- **データテーブル**: 6つのリソースを2列×3行で表示
- **通知システム**: 操作結果の成功・エラー通知

## 技術スタック

- **React 18**: UI フレームワーク
- **TypeScript**: 型安全性
- **Tailwind CSS**: スタイリング
- **Axios**: HTTP通信
- **React Hooks**: 状態管理

## 開発

### 新しいリソースの追加

1. `src/types/index.ts` に型定義を追加
2. `src/services/api.ts` にAPI関数を追加
3. `src/App.tsx` に状態管理とUIを追加

### スタイルのカスタマイズ

`src/App.css` でカスタムスタイルを追加できます。

## デプロイ

### 本番環境

```bash
npm run build
```

ビルドされたファイルをWebサーバーにデプロイしてください。

### Docker

```dockerfile
FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build
EXPOSE 3000
CMD ["npm", "start"]
```

## 注意事項

- このツールはデバッグ用です。本番環境では使用しないでください。
- データの削除は元に戻せません。重要なデータは事前にバックアップしてください。
- 依存関係のあるデータは適切な順序で作成してください。
