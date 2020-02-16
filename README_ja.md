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


## Getting Started
ダウンロードして使ってみたい場合、[GitHub - Releases]()からダウンロードしてください。

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

<!-- ### Installing
バイナリを任意の場所においてください。 -->
