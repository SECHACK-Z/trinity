# Trinity
個人開発者向けWebサービス運用フレームワーク

## What's this?
Webサービス運用時にある悩みを解決するためのフレームワークです。
個人開発者は、
- お金
- 時間
- 人手
がありません。

Trinityは
- ログ収集
- WebHookを用いた通知
- GitHubからのCD機能

を搭載しており、個人開発者の負荷を減らします。



### Prerequisites
ローカルでビルドする場合
- Node.js (v12.4以上)
- Go (v1.11以上)
    - [rakyll/tatik](https://github.com/rakyll/statik)
が必要です。

ビルドは`make`で完結します。
```
> make
```

## Requirements and Setup
[GitHub - Releases]()からダウンロードして試すことができます。

もしくは、手元でのビルドを試してください。

### Cloning the Repo
まずこのリポジトリをクローンしてください。

```sh
git clone git@github.com:sechack-z/trinity
```

ディレクトリに移動してください。
```sh
cd trinity
```

### Build as single binary
Trinityをかんたんに動かすために、シングルバイナリでのビルドをできるようにしました。

#### Install Dependencies.
- Node.js (above v12.4)
- Go (above v1.11)
- [rakyll/tatik](https://github.com/rakyll/statik)
    - シングルバイナリ化に使用   
- make

リポジトリのルートで`make`コマンドを使用してください。
```sh 
make build
```
もしくは、serverとclientそれぞれをビルドすることができます
### Client 
Clientは`client`ディレクトリに格納されています。
```
cd client
```
#### Node.js
Clientは[Node.js](https://nodejs.org/en/)で書かれています。
ビルドのために、`Node.js`をインストールしてください。
**Note**: `Node.js v12.4`で検証されていますが、いくつか問題があるかもしれません。何か問題が発生した場合、`v12.4`へのアップデートを検討してください。


#### Installing Dependencies
Yarnと呼ばれるパッケージマネージャを使用しています。
続けるにはYarnをインストールしてください。

依存パッケージをインストールしてください。
```sh
yarn
```

#### Building Client
Clientをビルドします。
```sh
yarn build:stage
```


### Server
Serverは`server`ディレクトリに格納されています。
#### Go
Serverは[Go](https://golang.org)で書かれています。
ビルドのために、`Go`をインストールしてください。
**Note**: `Go v1.11`で検証されていますが、いくつか問題があるかもしれません。何か問題が発生した場合、`v1.11`へのアップデートを検討してください。


#### Installing Dependencies
`Go modules`を使用しています。
そのため、ビルド時に依存パッケージは自動的にダウンロードされます。

#### Building Server
```sh
go build
```

