# Tages тестовое задание

## Билд 
```make build```

## Описание
После билда запустится сервер. Нужно запустить клиента ```go run pkg/client/client.go```  
Файлы юзера находятся в директории ```client_files```  
Сохраненные файлы на сервере находятся в директории ```server_saved```  
Запрошенные файлы от юзера находятся в директории ```client_requested_files```

## Пример работы 
Upload RPC && SendFiles RPC
```
2023-01-14T20:39:11.565+0300	info	client/client.go:38	Starting uploading: 20.png
[RESULT] success
2023-01-14T20:39:11.572+0300	info	client/client.go:38	Starting uploading: 40.png
[RESULT] success
2023-01-14T20:39:11.576+0300	info	client/client.go:38	Starting uploading: 65.png
[RESULT] success
2023-01-14T20:39:11.583+0300	info	client/client.go:38	Starting uploading: BATMAN_GOPHER.png
[RESULT] success
2023-01-14T20:39:11.584+0300	info	client/client.go:38	Starting uploading: BLUE_GLASSES_GOPHER.png
[RESULT] success
2023-01-14T20:39:11.588+0300	info	client/client.go:38	Starting uploading: GO_PARIS.png
[RESULT] success
2023-01-14T20:39:11.597+0300	info	client/client.go:46	Starting downloading files from server
2023-01-14T20:39:11.601+0300	info	client/client.go:65	Saving file: 20.png
2023-01-14T20:39:11.603+0300	info	client/client.go:65	Saving file: 40.png
2023-01-14T20:39:11.606+0300	info	client/client.go:65	Saving file: 65.png
2023-01-14T20:39:11.607+0300	info	client/client.go:65	Saving file: BATMAN_GOPHER.png
2023-01-14T20:39:11.608+0300	info	client/client.go:65	Saving file: BLUE_GLASSES_GOPHER.png
2023-01-14T20:39:11.611+0300	info	client/client.go:65	Saving file: GO_PARIS.png
2023-01-14T20:39:11.612+0300	info	client/client.go:70	Saving RPC was successful
ETC
ETC
ETC
```
GetAllFiles RPC  
```json
{
  "file": [
    {
      "name": "20.png",
      "createdAt": {
        "seconds": "1673717951",
        "nanos": 567212844
      },
      "updatedAt": {
        "seconds": "1673717951",
        "nanos": 567212947
      }
    },
    {
      "name": "40.png",
      "createdAt": {
        "seconds": "1673717951",
        "nanos": 573557538
      },
      "updatedAt": {
        "seconds": "1673717951",
        "nanos": 573557615
      }
    },
    {
      "name": "65.png",
      "createdAt": {
        "seconds": "1673717951",
        "nanos": 577868319
      },
      "updatedAt": {
        "seconds": "1673717951",
        "nanos": 577868403
      }
    },
    {
      "name": "BATMAN_GOPHER.png",
      "createdAt": {
        "seconds": "1673717951",
        "nanos": 583771660
      },
      "updatedAt": {
        "seconds": "1673717951",
        "nanos": 583771761
      }
    },
    {
      "name": "BLUE_GLASSES_GOPHER.png",
      "createdAt": {
        "seconds": "1673717951",
        "nanos": 585306124
      },
      "updatedAt": {
        "seconds": "1673717951",
        "nanos": 585306254
      }
    },
    {
      "name": "GO_PARIS.png",
      "createdAt": {
        "seconds": "1673717951",
        "nanos": 588444310
      },
      "updatedAt": {
        "seconds": "1673717951",
        "nanos": 588444429
      }
    }
  ]
}
```