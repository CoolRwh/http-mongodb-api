* 配置文件 `config/config.yml`

```yml
server:
  addr: ":8080"
mongodb:
  url: "mongodb://root:root@192.168.56.56:27017"
  minPoolSize: 10
  maxPoolSize: 100
```

> API 接口：

| 接口名称        | API                |
|-------------|--------------------|
| find        | api/v1/find        |
| fineOne     | api/v1/fineOne     | 
| installMany | api/v1/installMany |   
| installOne  | api/v1/installOne  |    
| updateById  | api/v1/updateById  |   
| updateOne   | api/v1/updateById  |    
| deleteById  | api/v1/deleteById  |   
| deleteMany  | api/v1/deleteMany  |  
