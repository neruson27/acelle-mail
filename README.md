## < Nombre pendiente >


# Divisiones
El proyecto se divide en 3 partes para su ejecución
1. *Server*: Inicia el servidor http
2. *Listener*: Inicia los watcher para los change streams de mongo db
3. *Jobs*: Inicia las tareas de sincronización de las companies y contacts con acelle mail

## Ejecución

### Make File (Linux)

Para instalar Make se puede ejecutar el siguiente comando en consola
```shell
sudo apt-get install make
```
Make utiliza un archivo llamado `Makefile` el cual dispone de un grupo de comandos a ejecutar, se ha implementado esto
para poder estar reduciendo tiempo de escritura/ejecución al momento de desarrollo.

Este es el listado de comandos que se disponen:

| Comando | Descripción |
|:---:| :----|
|build| genera el archivo binario del proyecto|
|clean| elimina si existe el archivo binario |
|all| Elimina, compila, ejecuta el proyecto mostrando la información de ayuda del mismo
|server| Elimina, compila, inicia el servidor http |
|listener | Elimina, compila, inicia el escuchador de cambios de streams de base de datos |
| jobs | Elimina, compila, inicia los jobs de integración de companies y de sincronización de contacts|

Uso:
```shell
make build # Unicamente genera el archivo binario del proyecto
make jobs # Elimina si existe el binario, compila el proyecto y ejecuta el proyecto desde los jobs
```

### Docker
Se tiene la creación de un contenedor de docker por medio de un `DockerFile` y al mismo tiempo la separación de los
diferentes divisiones del proyecto por medio de un `docker-compose.yml`.

Por medio del docker-compose se puede estar haciendo el build del contenedor y posterior levantar todos los servicios
o bien uno de ellos por separado.

Realizar el build del contenedor
```shell
docker-compose build
```

Iniciar todos los servicios
```shell
docker-compose up
```

Iniciar un unico servicio
```shell
docker-compose up server
```

Los servicios que se poseen son
1. server
1. listener
1. jobs

# Variables de Entorno
El proyecto tiene configuraciones con valores default, pero las mismas pueden estar siendo modificadas mediante
variables de entorno, esto se tienen para no estar almacenando archivo .env, sino que dependiendo el ambiente
se tienen diferentes configuraciones.

Se puede decir que se tienen diferentes variables de entorno por configuraciones. Dentro del archivo `config/keys.go` 
uno puede estar generando las variables de entorno disponibles, aunque los valores alli se tienen como por ejemplo 
`app.name` este valor como variable de entorno seria `APP_NAME`, el sistema lee en mayúsculas y con guion bajo. 

*Nota*
Sea por las configuraciones default o por variables de entorno, estos son valores que son utilizados por cualesquiera
de las 3 divisiones que posee el proyecto.

## Aplicación

| Key | Valor defecto | Descripción |
| :---: | :---: | :---|
| app.name| < Pendiente >| Nombre de la aplicación |
| secret | < ??? > | Secreto para estar haciendo la firma del jwt, en este escenario es para estar autenticando en el servidor http las request |

## Logger

| Key | Valor defecto | Descripción |
| :---: | :---: | :---|
| logger.code | zap | Indica que tipo de librería de logger se estará utilizando, actualmente se tiene implementada  `go.uber.org/zap`|
| logger.level | debug | Indica que tipo de nivel de logs se estarán registrando  (debug,info, warn, error)|

### Zap logger

| Key | Valor defecto | Descripción |
| :---: | :---: | :---|
| logger.zap.encoding| console | Indica el tipo de target es donde se estará utilizando el log (console, json) |
| logger.zap.development| true | Indica si la configuración actual es de desarrollo o producción |

Las siguientes configuraciones son keys que se utilizan para estar agregando al log que se genera, en el caso de que
no se quiera estar mostrando dicha parte del log se puede dejar la variable sin un valor.


| Key | Valor defecto | Descripción |
| :---: | :---: | :---|
| logger.zap.encoder.config.key.message | msg | Muestra el mensaje colocado en el logger |
| logger.zap.encoder.config.key.level | level | Muestra el tipo de logger que se genero (INFO, DEBUG, ERROR, FATAL) |
| logger.zap.encoder.config.key.time | ts | Muestra el momento en el que se genero el log |
| logger.zap.encoder.config.key.name | logger | (No tengo idea que es lo que muestra :D ) |
| logger.zap.encoder.config.key.name | <Vacio>  | Muestra el path del archivo que genero el logger |

## Mongo

| Key | Valor defecto | Descripción |
| :---: | :---: | :---|
| mongo.db.url | mongodb://localhost:27017/cliengo_dev_core | URL de conexión para la base de datos |
| mongo.db.name | cliengo_dev_core | es el nombre de la base de datos a la que se estará conectando, este valor es principalmente utilizado por el `listener`, pero en el caso que la URL no contenga el nombre de la DB es utilizado |
| mongo.db.timeout | 15 | Configuración en segundo de cuanto tiempo se estará tratando de hacer la conexión con la base de datos de forma exitosa |
| mongo.db.pool.size.min | 5 | Numero minimo de conexiones activas para el pool |
| mongo.db.pool.size.max | 15 | Numero maximo de conexiones activas para el pool |

## Acelle Mail

| Key | Valor defecto | Descripción |
| :---: | :---: | :---|
| acelle.mail.uri | https://emailmkt.stagecliengo.com | URL para estar comunicándose con el proyecto de acelle mail |
| acelle.mail.token | < ??? > | Token de autenticación para estar generando nuevas consumers en acelle mail |