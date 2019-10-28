## Siege 压力测试报告

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
| 1    | 10                  | 4.32               | 91.55 %      | 65                  | 3.92                 | 1.10                         |
| 2    | 100                 | 44.22              | 91.95 %      | 811                 | 3.25                 | 13.59                        |
| 3    | 250                 | 103                | 91.34 %      | 1856                | 3.31                 | 31.10                        |
| 4    | 500                 | 188.82             | 89.98 %      | 3332                | 3.36                 | 56.15                        |
| 5    | 1000                | 375.04             | 90.36 %      | 6909                | 3.26                 | 115.11                       |

#### 结果分析：

* 在GET index页面时，直接在前端利用CDN加载React库，耗费了大量的时间
* 长短链接的储存都在内存中进行，效率较低
* 因为Siege的最大并发数为1000，无法进行1000以上的测试



### Stage 2

| ID   | Setting Concurrency | Actual Concurrency | Availability | Transactions (hits) | Response time (secs) | Transaction rate (trans/sec) |
| ---- | ------------------- | ------------------ | ------------ | ------------------- | -------------------- | ---------------------------- |
| 1    | 10                  |                    |              |                     |                      |                              |
| 2    | 100                 |                    |              |                     |                      |                              |
| 3    | 250                 |                    |              |                     |                      |                              |
| 4    | 500                 |                    |              |                     |                      |                              |
| 5    | 1000                | 380.9              | 90.59%       | 6868                | 3.34                 | 114.2                        |

#### 结果分析

* 因为服务器性能有限，在链接MySQL数据库时抛出了`Too many connections`错误，此处为服务器性能的瓶颈。

