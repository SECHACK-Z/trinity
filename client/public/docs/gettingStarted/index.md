# Getting Started
このページは<かっこいい名前>の起動に関する資料です。

# <かっこいい名前>を試す
<かっこいい名前>を試すには、以下の手順に従ってください。
1. <かっこいい名前>のビルド
2. `config.yaml` の編集
3. 実行

## <かっこいい名前>のビルド
<かっこいい名前>は大きく分けて２つに別れています。
```
.
├── Dockerfile
├── Makefile
├── client
├── server
└── test
```
- `client/` がWebUIに関してのディレクトリ
- `server/` が<かっこいい名前>本体のディレクトリ

<かっこいい名前>を使用するためには、ルートディレクトリで`make`コマンドを使用し、ビルドを行ってください。
```
$ make
```
ここでエラーが出る場合、ビルドに使用しているコマンドがシステムに存在しないことが原因であることが多いです。
`npm`, `go`コマンドが正しくインストールされていることを確認してください。


ビルドが完了すると、server以下に`app`というファイルが作成されています。
まずはそのファイルを実行してみましょう。
```
$ cd server
$ ./app
2020/01/12 22:36:43 new settings apply
2020/01/12 22:36:43 director register
2020/01/12 22:36:43 webUI start
⇨ http server started on [::]:8080
```

なにもエラーが出なければ正しくビルドが完了しています。

## <かっこいい名前>の実行
<かっこいい名前>を正しく動作させるためには、`config.yaml`を作成する必要があります。

以下に設定項目と、設定ファイルの例を示します。
|設定項目|例|説明|
|---|---|---|
|proxy|localhost:8080|プロキシ先|
|host|a.sechack.org|割り当てたいホスト名|
|https|true|https通信をするか|
|forcehttps|false|https通信を強制するか|
|default|true|一致するホスト名がなかったときの転送先にするか|
|healthcheck| false| ヘルスチェック機能を有効にするか|
|repository|https://github.com/onsd/misc|CD機能を有効にするか|

設定ファイルの例
```yaml
targets
-   proxy: localhost:9091
    host: a.sechack.org
    https: false
    forcehttps: false
    default: true
    healthcheck: false
    repository: https://github.com/onsd/misc
-   proxy: localhost:9092
    host: b.sechack.org
    https: false
    forcehttps: false
    default: false
    healthcheck: false
```

### ヘルスチェック機能
[リンク](http://localhost:8080/#/documentation/index)
### CD機能
[リンク](http://localhost:8080/#/documentation/index)
