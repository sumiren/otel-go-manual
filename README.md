## Concept

* GoにはOpenTelemetryの自動計装がない（2024年3月時点、公式）。手動計装するとどんなものかと思い試した
  * Gin
  * SQLクライアント
  * HTPクライアント

## 結果


* JVMのようなエージェントの仕組みやNode.jsのようにtracing.jsをアプリのコードベースと別に用意したりはできない
* SQLやHTTPのアウトバウンド通信など主要なモジュールにラッパーが用意されているため、例えば`sql.Open`を`otelsql.Open`のように書き換えていく
  * [Gin](https://github.com/sumiren/otel-go-manual/blob/main/main.go#L62C2-L62C38)
  * SQL
    * [connect](https://github.com/sumiren/otel-go-manual/blob/main/main.go#L64)
    * [query](https://github.com/sumiren/otel-go-manual/blob/main/db.go#L19)
      * 直前の`tracer.start`は別のカスタムスパンなので関係ない
  * [HTTPクライアント](https://github.com/sumiren/otel-go-manual/blob/main/http.go#L20) 
    * 同様に直前の`tracer.start`は関係ない
* 全体的にGinのリクエストコンテキストを引き回す必要があるのがしんどい

これで取れるトレースのイメージ

![image](img.png)

## 感想

* ラッパーのおかげでtraceにattributeを詰める処理などを書く必要はない点はよかった
* コードベース次第で導入コストが大きく変わってくるという印象
  * 色んな場所で`sql.Open`とかしてるとしんどい
  * コンテキストの引き回しされてないと大規模な改修になる