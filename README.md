# 使用说明
thz-spider.exe  redis-host:port redis-passwd [启动类型 master/slave]

**先启动 master * 1 再启动 slave * N 重复启动mater会影响任务队列**

默认输出到工作目录下out文件夹

基于colly开发
> http://go-colly.org/

目前基于Redis共享任务队列 可分布式部署
