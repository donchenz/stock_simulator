
目录结构
==============
middleware       	            <em>//消息处理中间件</em>  
&ensp;&ensp; connection.go  	<em>//处理和客户端的连接</em>  
&ensp;&ensp; hub.go 			    <em>//管理连接的集合，可以注册、注销连接，广播消息到各条连接等</em>  
&ensp;&ensp; main.go 			    <em>//程序入口</em>  
&ensp;&ensp; redisclient.go 	<em>//处理redis的存取</em>  
web                 <em>//客户端</em>  
&ensp;&ensp; client.html  
stock_gen.py	      <em>股票⾏情模拟⽣成器</em>  


安装 & 运行
==============

golang用到的第三方库:
-------------------
github.com/gorilla/websocket  
github.com/go-redis/redis

运行:
-------------------
1. 项目需要用到redis，默认端口6379;
2. 运行中间件。进入middleware目录，执行 go run *.go (或编译之后运行);
3. 运行股票数据生成程序 python stock_gen.py;
4. 打开浏览器，输入localhost:8080 进入客户端.


设计动机 & 考虑因素
==============
1. 股票⾏情模拟⽣成器和消息推送中间件之间的连接有可能丢数据，使用消息队列解决;
2. 客户端和服务器之间需要心跳机制;
3. 客户端可能有成千上万个，单进程的服务器不能支撑，需要考虑分布式;
4. 需要考虑消息在特定时候的浪涌情况, 因此评估系统服务能力时应按峰值评估。



技术实现
==============
1. stock_gen随机生成股票数据，发送给消息处理中间件middleware。stock_gen和middleware
   之间的通信通过消息队列进行。demo通过redis的lpush/lpop实现了简易的消息队列。stock_gen
   不断把数据推送至redis，middleware不断从redis获取数据。redis定时将数据同步到db，避免
   因redis crash导致的丢消息。
2. middleware通过websocket的方式不断推送消息到客户端。每条连接产生一条connection。
   connection内置一个发送channel。
3. middleware从redis取出一条数据后，放到hub的broadcast channel。hub 从broadcast channel
   取出消息，对每一条在线的连接进行广播。
4. 客户端定时向middleware发送心跳包，如果心跳超时，middleware会清理掉超时的连接。


扩展
==============
1. 单进程不足以支持成千上万个客户端，考虑通过多个middleware提供服务。redis消息队列设计成环形队列
   模式，即队列达到最大长度后，每push一条消息进队头，就要从队尾pop一条消息扔掉。每个middleware
   不是通过pop消费消息，而是通过指针记录已经消费到哪条消息。这样避免每个middleware都要保存一条
   消息队列。
2. 扩展消息协议的格式，以便支持更复杂的协议。
3. 客户端可能需要订阅某些股票的行情，因此考虑增加订阅功能，只给客户端推送订阅的数据。


其他
===============
⽬前的实现⽅案，有⼀个⼩问题，就是当⼀个新客户端连接上中间件后，可能短时间
内没有数据（⽐如碰巧下⼀轮⾏情还没到，⽽且到的也可能只是其中⼀两只股票），
所以我们也希望客户端连接上来后，先获得此刻之前的历史数据。你会考虑怎么解决
这个问题？

 ---- 我会考虑把历史数据存到服务器中。比如，消息中间件起一个定时任务，定期
      从其他地方获取历史数据，存到db。当客户端的连接建立后，首先发送获取历史
      数据的请求，然后服务端返回历史数据，客户端收到历史数据后，再发送开始推
      送的请求给服务端，然后服务端才开始推送。获取历史消息和推送不能并行，否则
      有可能先显示了推送消息，又被历史消息覆盖掉。










