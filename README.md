### Клиент для https://vk.com
#### Установка
```
go get github.com/FireGM/vk
```
### Использование
Создание клиента через токен
```
client.DefaultClient("YOUR_TOKEN")
```
Создание клиента через логин/пароль
```
GetClientWithLoginPassword(username, password)
```
Запросы
```
package main

import (
	"github.com/FireGM/vk/client"
	"net/url"
	"fmt"
)

func main() {
	cl := client.DefaultClient("")
	values := url.Values{}
	values.Add("group_id", "1")
	req, _ := cl.MakeRequest("groups.getById", values)
	res, _ := cl.DoBytes(req)
	fmt.Println(string(res))
	//{"response":[{"id":1,"name":"ВКонтакте API","screen_name":"apiclub","is_closed":0,"type":"group"}]}
}
```

