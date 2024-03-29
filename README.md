# th-bingo-tools

![](https://img.shields.io/github/go-mod/go-version/CuteReimu/th-bingo-tools "语言")
[![](https://img.shields.io/github/actions/workflow/status/CuteReimu/th-bingo-tools/golangci-lint.yml?branch=master)](https://github.com/CuteReimu/th-bingo-tools/actions/workflows/golangci-lint.yml "代码分析")
[![](https://img.shields.io/github/contributors/CuteReimu/th-bingo-tools)](https://github.com/CuteReimu/th-bingo-tools/graphs/contributors "贡献者")
[![](https://img.shields.io/github/license/CuteReimu/th-bingo-tools)](https://github.com/CuteReimu/th-bingo-tools/blob/master/LICENSE "许可协议")

请配合 [th-bingo服务端](https://github.com/CuteReimu/th-bingo) 和 [th-bingo客户端](https://github.com/Death-alter/th-bingo) 食用

**目前只有Windows版本**

## 编译

```shell
go build -o th-bingo-tools.exe
```

## 运行

运行后监听 9760 端口，访问 ws://127.0.0.1:9760/ws 即可获取选卡、收卡的回调。

`test.exe -h` 可以查看更多命令行参数的用法

回调就是一个json，如下所示：

```json
{
  "game": 18,
  "id": 1,
  "event": 1,
  "mode": 1,
  "role": "Reimu",
  "rank": "L",
  "score": 123450
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
| score | int | 符卡最高分数，一般只有Spell Practice才有，没有分数则没有此字段     | 

角色的英文名：（对于有不同子机的作品，在英文名后加字母表示，例如：ReimuA、SakuyaB。特别地，天空璋不区分子机，鬼形兽用WOE三个字母表示三个支援。）

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

## 开发相关

总体思路就是用CE找一下存放符卡history数据的内存，然后找到规律。

下面用th10举个例子：

https://github.com/CuteReimu/th-bingo-tools/blob/b1fb048d9386b1d5f6eb1f46a533ec7b3d68c88f/listener_th10.go#L32-L35

上述代码的意思是：从进程的基址开始，往后`0x7783C`个字节的位置，存放了指向角色数据结构体的指针，我们利用这个指针找到这个结构体后，再向后移动`20`个字节的位置，就是第一个角色（灵梦A）的id；而往后`0x5A4`个字节的位置，就是第一个角色的符卡数据的结构体的起始地址。每个角色的数据距离上一个角色的数据`0x437C`个字节的距离。用CE找到这些地址之后，我们把内容读出来即可。
