# csv-join

按共享列把两份 CSV 做 inner / left join（CLI）。

## 用法

```
csv-join -left a.csv -right b.csv [-key id] [-how inner] [-out -]
```

- `-left`   左表路径（必填）
- `-right`  右表路径（必填）
- `-key`    两边都有的连接列，默认 `id`
- `-how`    `inner`（默认）或 `left`
- `-out`    输出路径，`-` 表示标准输出

左表带 UTF-8 BOM 时会剥掉；以 `#` 开头的行会跳过。数据行列数必须和表头一致。
