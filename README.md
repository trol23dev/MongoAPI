
## Author

- [Ramón L.](https://www.github.com/trol23dev)


# API for MongoDB

The API developed for this project serves as the backbone of the system, facilitating communication with a MongoDB database. Here's a brief description of the API:

The API is built using Express.js, a fast, unopinionated, and minimalist web framework for Node.js. It provides a robust set of features for building web applications and APIs, including routing, middleware support, and HTTP utilities.

The API exposes several endpoints that handle various operations related to tickets and inventory management. These endpoints allow to interact with the database seamlessly.

The API leverages MongoDB as the database management system, providing a flexible and scalable solution for storing and retrieving data. It utilizes the MongoDB Node.js driver to establish connections with the database and execute queries to perform CRUD (Create, Read, Update, Delete) operations on documents.

Overall, the API acts as a bridge between a project and the database, allowing for seamless communication.
## Deployment

To deploy this project run

```bash
  npm run build
```

```bash
  npm run start
```
## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`MONGO_URI`


## API Reference

#### Get all items

```http
  GET /api/:database/:collection
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `database` | `string` | **Required**. Your database name. |
| `collection` | `string` | **Required**. The collection you want to get the information. |

#### Post item

```http
  POST /api/:database/:collection
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `database` | `string` | **Required**. Your database name. |
| `collection` | `string` | **Required**. The collection you want to push the information. |

#### Delete item

```http
  DELETE /api/:database/:collection/:id
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `database` | `string` | **Required**. Your database name. |
| `collection` | `string` | **Required**. The collection you want to delete the information. |
| `collection` | `string` | **Required**. The id of the information you want to delete. |

![Logo](https://i.postimg.cc/WbyjnZyT/Logotipo-Birdy-Project-03.png)

