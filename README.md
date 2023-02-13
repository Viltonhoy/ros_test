Описание выполняемого тестового задания.
=========================================

### Первое и второе задание:

## Описание веб-приложения:

```golang
func main() {
	http.HandleFunc("/", helloDocker)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func helloDocker(w http.ResponseWriter, r *http.Request) {
	var s = "Hello Docker!"

	if ev := os.Getenv("VALUE"); ev != "" {
		s = ev
	}
	fmt.Fprint(w, s)
}
```

- В функции main описывается хендлер, в функции helloDocker – стандартный вывод строки “ Hello Docker!” и использование переменной окружения, для замены данной строки.

- Наше приложение выводит на локальном хосте с указанным портом сообщение “Hello Docker!”

## Cоздание и описание докер файла:

```dockerfile
# syntax=docker/dockerfile:1
ARG GO_VERSION=1.17

FROM golang:${GO_VERSION}-alpine AS builder
ENV GO111MODULE=on
WORKDIR /src      
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build \
 -installsuffix `static` \
 -v -o /app cmd/server/main.go

FROM gcr.io/distroless/static-debian11
WORKDIR /
COPY --from=builder /app /app
EXPOSE 8080 8080
USER nobody:nobody
ENTRYPOINT ["/app"]
```

1. Поле ARG – задает переменную для передачи Docker во время сборки, а именно – версию образа Golang. 
2. FROM — задаёт базовый (родительский) образ. В данном случае базой этого образа идет официальный образ Golang c тегом {*указанная версия в аргументе*}-alpine.
3. ENV — устанавливает постоянные переменные среды. В данном приложении используется переменная окружения для указания новой строки при выводе.
4. WORKDIR — задаёт рабочую директорию для следующей инструкции.
5. СOPY — копирует в контейнер файлы и папки. Сообщаем Docker о том, что нужно взять папки из локального контекста сборки и добавить их в текущую рабочую директорию образа.
6. RUN — выполняет команду и создаёт слой образа. Используется для установки в контейнер пакетов, а именно – наше веб-приложение.
7. EXPOSE — указывает на необходимость открыть порт 8080
8. USER - устанавливает имя пользователя или UID для использования при запуске образа и для любых инструкций RUN, CMD и ENTRYPOINT, которые следуют за ним в Dockerfile. В нашем случае имя не указано. 

После описания Docker файла, с помощью консольной команды, билдим наш образ.
```
docker build -f 'build/package/Dockerfile' -t app .
```  

Запускаем Docker контейнер

```
docker run -it -p 8080:8080 app
```

Результат:

![](https://github.com/Viltonhoy/ros_test/blob/master/images/a.png)

Для изменения строки запускаем новый контейнер командой с описанием переменной окружения.
```
docker run -it -p 8080:8080 -e "VALUE= Hello Sibintek!" app
```

Результат:

![](https://github.com/Viltonhoy/ros_test/blob/master/images/b.png)

### Hазвертываниt приложения в Kubernetes

- Для запуска одноузлового кластера Kubernetes на виртуальной машине используется гипервизор Hyper-v.

Запускаем локальный кластер Kubernetes командой
```
minikube start --vm-driver=hyperv
```

## Описание приложения с помощью Kubernetes YAML

Создаем файл развертывания в корневом каталоге с именем kuber-pod.yaml. Этот файл развернет приложение в движке Kubernetes. В метадате указываем лейбл пара ключ: значение, для дальнейшего использования в селекторе сервиса.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: app-kuber
  labels:
    app: web 
spec:
  containers:
  - name: sibintek-app
    image: app
    imagePullPolicy: IfNotPresent
    ports:
    - containerPort: 8080
      protocol: TCP
```

Создаем файл сервиса c именем my-service.yaml, для абстрагирования доступа к модулю Kubernetes, используя в селектору пару ключ:значение.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-service
spec:
  type: NodePort
  selector:
    app: web 
  ports:
    - protocol: TCP
      port: 8080
      targetPort: 8080
      nodePort: 30001 
```

Так как работа ведется локально, загружаем локальный докер-образ косандой.

```
minikube image load app
```

Разворачиваем наши экземпляры cервиса и приложения в движке Kubernetes командами.

```
kubectl apply -f kuber-pod.yaml
```

```
kubectl apply -f my-service.yaml
```

Проверяем работоспособность пода командой
```
kubectl get pods
```

