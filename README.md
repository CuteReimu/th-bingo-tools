# th-bingo-tools

请配合 [th-bingo服务端](https://github.com/CuteReimu/th-bingo) 和 [th-bingo客户端](https://github.com/Death-alter/th-bingo) 食用

**目前只有Windows版本**

## 编译

```shell
go build -o th-bingo-tools.exe
```

## 运行

运行后监听 9961 端口，访问 ws://127.0.0.1:9961/ws 即可获取选卡、收卡的回调

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
| role  | str | 角色的英文名（见下方表格）                              |
| rank  | str | 难度（E、N、H、L、EX、PH）                          | 

角色的英文名：（对于有不同子机的作品，在英文名后加字母表示，例如：ReimuA、SakuyaB。特别地，天空璋仍用ABCD表示四个子机，鬼形兽用WOE三个字母表示三个支援。）

| 英文名     | 含义   | 备注      |
|---------|------|---------|
| Reimu   | 灵梦   |         |
| Marisa  | 魔理沙  |         |
| Sakuya  | 咲夜   |         |
| Sanae   | 早苗   |         |
| Youmu   | 妖梦   |         |
| RY      | 结界组  | 仅在永夜抄中有 |
| MA      | 咏唱组  | 仅在永夜抄中有 |
| SR      | 红魔组  | 仅在永夜抄中有 |
| YY      | 幽冥组  | 仅在永夜抄中有 |
| Yukari  | 紫    | 仅在永夜抄中有 |
| Alice   | 爱丽丝  | 仅在永夜抄中有 |
| Remilia | 蕾米莉亚 | 仅在永夜抄中有 |
| Yuyuko  | 幽幽子  | 仅在永夜抄中有 |
| Reisen  | 铃仙   | 仅在绀珠传中有 |
| Cirno   | 琪露诺  | 仅在天空璋中有 |
| Aya     | 射命丸文 | 仅在天空璋中有 |
