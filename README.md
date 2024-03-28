
# Server

## Descripción del Proyecto

El objetivo de este proyecto es desarrollar un servidor en el lenguaje de programación Go que actúe como intermediario entre una base de datos PostgreSQL y un programa en Vue.js. El servidor proporcionará una API para permitir el intercambio de datos entre la base de datos y el programa en Vue.js.

### Requisitos Previos

- Instalación de Go: [Descargar e instalar Go](https://golang.org/doc/install)

- Instalación de Vue.js: [Guía de inicio de Vue.js](https://vuejs.org/v2/guide/installation.html)

### Configuración de la Base de Datos PostgreSQL

Asegúrese de tener una base de datos PostgreSQL en funcionamiento. Puede utilizar herramientas como pgAdmin para administrar su base de datos.

### Configuración del Proyecto

1. Clone este repositorio en su máquina local:

   ```bash
   git clone https://github.com/jerapepe/server.git
   ```

2. Navegue al directorio del proyecto:

   ```bash
   cd server
   ```

3. Instale las dependencias del proyecto:

   ```bash
   go mod tidy
   ```


## Uso

Para ejecutar el servidor, simplemente ejecute el siguiente comando desde el directorio raíz del proyecto:

```bash
go run main.go
```

El servidor comenzará a escuchar en el puerto especificado.

## Licencia

Este proyecto está bajo la licencia [MIT License](LICENSE).

---
