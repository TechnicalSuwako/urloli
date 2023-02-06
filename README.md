# URLロリ
クッソ小さいURL短縮作成ソフトだわ〜♡

## 使い方
```sh
cp links.sample.json links.json
nvim links.json
```

nginxのコンフィグで：
```
  location / {
    add_header Permissions-Policy interest-cohort=();
    rewrite ^/(.+)$ /index.php?url=$1 last;
    try_files $uri $uri/ /404.html;
  }
```

links.jsonファイルの中に：
```
{
  "hogehoge": "https://076.moe"
}
```

https://（ドメイン名）/hogehoge にアクセスすると、https://076.moe に移転されます。

## APIの使い方

### 短縮URLの創作

```
METHOD: POST
URL: https://urlo.li
PARAM: {
  sosin: 1
  api: 1
  newurl: "(元URL)"
}

OUT: { "res": "(5英文字)" }
```

### 短縮したURLの確認

```
METHOD: GET
URL: https://urlo.li/kk1v9?api=1
OUT: { "res": "https://076.moe" }
```
