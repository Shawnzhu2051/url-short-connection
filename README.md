# URL 短连接系统

#### Introduction

用Go实现一个URL短连接系统，类似于 http://dwz.cn/ ，采用 https://github.com/gin-gonic/gin 这个框架

1. 第1阶段可以先在单进程实现，长短url实现在存放在进程内就好
2. 第2阶段可以将长短url的关联关系放到MySQL中，MySQL驱动采用https://github.com/jinzhu/gorm
3. 第3阶段为了提高性能，将热点url放到Redis中缓存，全量数据仍然放到MySQL中 https://github.com/go-redis/redis
4. 每一阶段实现完，做压测、性能调优。没有实验环境的可以尝试使用阿里云

#### Tips

1. 因为go-redis目前官方只支持使用go module进行管理，如果有混用govendor的话，可能会导致本地自定义的`utils/URLShortener`包无法使用。可以在`go.mod`中添加`replace utils/URLShortener v0.0.0 => /root/go/src/github.com/shawnzhu/url-short-connection/utils/URLShortener`来获取本地自定义包。



## Demo

Link: [http://45.77.114.214:8080/](http://45.77.114.214:8080/)

![](./image/demo.png)

## Siege 压力测试

使用Siege来进行压力测试，测试指令：

```bash
siege -c [concurrency number] -t 1M -i -f urlTest.txt
```

`-c` 指定并发数

`-t` 指定运行时间1分钟 

`-i` 随机抽取url

`-f` 使用urlTest.txt文件

Vulter服务器性能如下：

![](./image/quality.png)

### Stage 1

| ID   | Setting Concurrency | Actual Concurrency | Availability | Transactions (hits) | Response time (secs) | Transaction rate (trans/sec) |
| ---- | ------------------- | ------------------ | ------------ | ------------------- | -------------------- | ---------------------------- |
| 1    | 10                  | 9.95               | 100 %        | 855                 | 0.69                 | 14.32                        |
| 2    | 100                 | 81.96              | 99.86 %      | 6610                | 0.74                 | 110.41                       |
| 3    | 250                 | 209.17             | 99.80 %      | 14960               | 0.84                 | 250.33                       |
| 4    | 500                 | 286.33             | 98.35 %      | 17334               | 0.98                 | 293.40                       |
| 5    | 1000                | 460.16             | 96.89 %      | 20340               | 1.35                 | 340.53                       |

### Stage 2

| ID   | Setting Concurrency | Actual Concurrency | Availability | Transactions (hits) | Response time (secs) | Transaction rate (trans/sec) |
| ---- | ------------------- | ------------------ | ------------ | ------------------- | -------------------- | ---------------------------- |
| 1    | 10                  | 9.95               | 100.00 %     | 842                 | 0.70                 | 14.24                        |
| 2    | 100                 | 80.95              | 99.70 %      | 3975                | 1.22                 | 66.25                        |
| 3    | 250                 | 181.17             | 98.91 %      | 10639               | 1.02                 | 177.82                       |
| 4    | 500                 | 266.25             | 97.69 %      | 12653               | 1.26                 | 210.78                       |
| 5    | 1000                | 392.45             | 96.07 %      | 17044               | 1.38                 | 283.83                       |

### Stage 3

| ID   | Setting Concurrency | Actual Concurrency | Availability | Transactions (hits) | Response time (secs) | Transaction rate (trans/sec) |
| ---- | ------------------- | ------------------ | ------------ | ------------------- | -------------------- | ---------------------------- |
| 1    | 10                  | 9.53               | 100 %        | 734                 | 0.77                 | 12.32                        |
| 2    | 100                 | 90.40              | 99.96 %      | 7288                | 0.73                 | 123.36                       |
| 3    | 250                 | 172.20             | 99.07 %      | 10739               | 0.96                 | 178.95                       |
| 4    | 500                 | 277.74             | 98.32 %      | 17771               | 0.93                 | 298.02                       |
| 5    | 1000                | 494.99             | 97.09 %      | 20463               | 1.45                 | 341.51                       |

#### 结果分析

* 在GET index页面时，直接在前端利用CDN加载React库，耗费了大量的时间
* 因为Siege的最大并发数为1000，无法进行1000以上的测试
* 在分别使用MySQL和redis后，性能较之前有一定提升。但效果并不特别显著。猜测是因为1000的并发量过小，没够体现出MySQL和redis的优点。
* Siege的测试结果不稳定，反复测试时结果相差较大，此处取的多次测试中最好的结果。猜测是因为Vulter服务器位于境外，网络不稳定造成的。

