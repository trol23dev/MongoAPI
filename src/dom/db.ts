import { MongoClient, Db } from 'mongodb';
import Debug  from './debug';

export default class DB {
    private client: MongoClient;
    private db: Db;

    constructor(databaseName: string) {
        this.client = new MongoClient(`${process.env.MONGO_URI}`);
        this.db = this.client.db(databaseName);
    }

    async connect() {
        try {
            await this.client.connect();
            new Debug('success').print("🍃 MongoDB connection has been established.")
        } catch (error) {
            new Debug('error').print("🍃 MongoDB connection has been failed.")
        }
    }

    async close() {
        await this.client.close();
        new Debug('system').print("🍃 MongoDB connection has been closed.")
    }

    getDatabase() {
        return this.db;
    }
}