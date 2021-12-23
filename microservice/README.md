# TODO List
------------------------------------------------------------------------------------------
### _架構面_
- [ ] strategy-manager-mongo 與 openfaas function 的 synchronization
  
### _Grpc faas service_

- [x] 需要新增 update strategy 功能 
- [x] 用 crontab 更新 strategy 的排程
- [x] ~~deploy function 需要 pod request/limit 
      --> 可以用limitRange + resoureQuota去解（？~~
- [ ] deploy function前要先到 user resource service 問可調度的cpu/memory     
- [ ] 評估是否需要secure connection for demo

### _user resource service_

- [x] resourceQuota & limitRange 在 create namespace時一併執行
- [ ] 讓admin調整使用者的resources
- [ ] scratch idea: 用cognito創admin scope?

### _strategy manager service_

- [x] 提供 crontab 更新function的route
- [x] 從這邊發出 DELETE request，要找辦法讓 grpc faas知道 delete request ~~(暫定kafka)~~ 
    ------> *後來使用 grpc 直接呼叫* 
- [ ] 解析bearer token，如果呼叫的path跟uid衝突，要處理
- [ ] kafka consumer 跟 handler 應該透過 channel 溝通
- [x] feature request: (1) 調整strategy的resources (2)label annotations for schedule
- [x] feature request: 應該要讓使用者隨意調req/limit, 100m/128Mi & 750m/512Mi 只是預設值
    
### _front end_
- [x] 呈現目前使用者所有strategy的頁面
- [ ] delete 跟 trigger strategy
- [ ] show logs of strategy
- [ ] 調整strategy的cpu/memory

### _istio_
- [x] 需修改configMap，user strategy pod 不該有sidecar-proxy
  ----> 安裝的時候直接edit configMap
- [x] envoyfilter用來限制C9用戶的封包發送
- [x] backend service 新增 requestAuthenticaion / AuthorizationPolicy
- [ ] service port 用 appProtocol 的差別？
- [x] virtual service 為何跳不到 grpc service? --> istio的vs host還是要用原本K8s的FQDN
