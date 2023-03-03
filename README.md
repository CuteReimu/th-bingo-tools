# th-bingo-tools

请配合 [th-bingo服务端](https://github.com/CuteReimu/th-bingo) 和 [th-bingo客户端](https://github.com/Death-alter/th-bingo) 食用

**目前只有Windows版本**

## 编译

```shell
go build -o th-bingo-tools.exe
```

## 运行

~~运行后监听 9961 端口，访问 ws://127.0.0.1:9961/ws 即可获取选卡、收卡的回调~~ （还没实现）

回调就是一个json，如下所示：

```json
{
  "game": 18,
  "id": 1,
  "event": 1,
  "mode": 1,
  "role": "Reimu",
  "rank": "L"
}
```

各字段含义如下：

| 字段    | 类型  | 含义                                         |
|-------|-----|--------------------------------------------|
| game  | int | 作品号，6到18，目前只支持18                           |
| id    | int | 符卡id，游戏里可以看到 No.xx                         |
| name  | str | 符卡名，不一定有的字段，也不一定是中文，不建议使用                  |
| event | int | 事件，0-进入符卡，1-收取符卡                           |
| mode  | int | 0-游戏模式或Practice Start<br/>1-Spell Practice |
| role  | str | 角色的英文名                                     |
| rank  | str | 难度（E、N、H、L、EX、PH）                          | 