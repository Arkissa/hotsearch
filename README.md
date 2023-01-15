# HOT SEARCH

这是一个热搜程序，自动筛选关键词


## 参数
```shell
$ ./hotsearch -h
  -d string
    	The Database file Path (default "hotsearch.db")
  -h	Help
  -k string
    	The keyword file Path (default "keywords.csv")
```

## Q & A

1. 程序多久执行一次
    > 默认半小时执行一次

2. 如何存储的数据
    > 使用的是sqlite3数据库会自动创建数据库文件并且自动创建表

3. 关键词的逻辑是什么
    > 关键词判断逻辑仅仅是判断标题是否出现关键词，并且是热加载的关键词文件
