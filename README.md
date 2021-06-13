# 作业

基于 errgroup 实现一个 http server 的启动和关闭，以及 linux signal 信号的注册和处理，要保证能够 一个退出，全部注销退出。



# 细节

http server 与 信号的注册和处理都会阻塞，可以使用group.Go启动goroutine处理

http server 启动一个done的chan做监听处理，开放close接口做关闭http server 的控制

linux signal 信号本身是chan阻塞，有信号就可以解除
