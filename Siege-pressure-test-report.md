## Siege 压力测试报告

使用Siege来进行压力测试，测试指令：

```bash
siege -c 1000 -t 1M -i -b -f urlTest.txt
```



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
| 1    |                     |                    |              |                     |                      |                              |
| 2    |                     |                    |              |                     |                      |                              |
| 3    |                     |                    |              |                     |                      |                              |
| 4    |                     |                    |              |                     |                      |                              |
| 5    | 1000                | 380.9              | 90.59%       | 6868                | 3.34                 | 114.2                        |

Too many connections