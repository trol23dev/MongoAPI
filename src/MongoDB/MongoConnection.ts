import { MongoClient, Db } from "mongodb";
import { color } from "../functions"; // Importamos nuestra función para los colores

export default class MongoDB {
    private client: MongoClient;
    private db: Db;
    
    constructor(databaseName: string) {
        // Creamos una instancia de MongoClient utilizando la URL de conexión en la variable de entorno MONGO_URI
        this.client = new MongoClient(`${process.env.MONGO_URI}`);
        // Obtenemos una referencia a la base de datos especificada en la variable de entorno MONGO_DATABASE_NAME
        this.db = this.client.db(databaseName);
    }

    async connect() {
        try {
            // Intentamos establecer la conexión con MongoDB
            await this.client.connect();
            // Si la conexión es exitosa, imprimimos un mensaje de registro en la consola
            console.log(`🍃 MongoDB connection has been ${color("variable", "established.")}`);
        } catch (error) {
            // Si la conexión falla, imprimimos un mensaje de error en la consola
            console.log(`🍃 MongoDB connection has been ${color("error", "failed.")}`);
        }
    }

    async close() {
        // Cerramos la conexión con MongoDB
        await this.client.close();
        // Imprimimos un mensaje de registro en la consola indicando que la conexión se ha cerrado
        console.log(`🍃 MongoDB connection has been ${color("error", "closed.")}`);
    }

    getDatabase() {
        // Devolvemos una referencia a la base de datos
        return this.db;
    }
}
