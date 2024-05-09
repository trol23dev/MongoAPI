import DB from "../dom/DB";
import Debug from "../dom/debug";

import express, { Request, Response } from "express";
import bodyParser from "body-parser";

import { ObjectId } from "mongodb";

export default class API {
    private port: number;
    private app: express.Application;

    constructor(portNumber?: number) {
        this.port = portNumber || 3000;
        this.app = express();
    }

    async start() {
        this.app.use(bodyParser.json());

        this.app.get(`/api/:database/:collection`, async (req: Request, res: Response) => {
            const db = new DB(req.params.database);
            await db.connect();

            try {
                const collection = db.getDatabase().collection(req.params.collection);
                const documents = await collection.find({}).toArray();
                res.json(documents);
            } catch (e) {
                res.status(500).json({ error: "Error al obtener documentos" });
            } finally {
                await db.close();
            }
        });

        this.app.post(`/api/:database/:collection`, async (req: Request, res: Response) => {
            const db = new DB(req.params.database);
            await db.connect();

            try {
                const collection = db.getDatabase().collection(req.params.collection);
                const newDocument = req.body;

                const result = await collection.insertOne(newDocument);

                if (result) {
                    res.status(200).json({
                        message: "Data insertada correctamente"
                    });
                } else {
                    res.status(500).json({ error: "Error al crear el documento" });
                }
            } catch (e) {
                res.status(500).json({ error: "Error al crear el documento" });
            } finally {
                await db.close();
            }
        });

        this.app.delete(`/api/:database/:collection/:id`, async (req: Request, res: Response) => {
            const db = new DB(req.params.database);
            await db.connect();

            try {
                const collection = db.getDatabase().collection(req.params.collection);
                const result = await collection.deleteOne({ _id: new ObjectId(req.params.id) });

                if (result) {
                    res.status(200).json({
                        message: "Data eliminada correctamente"
                    });
                } else {
                    res.status(500).json({ error: "Error al eliminar el documento" });
                }
            } catch (e) {
                res.status(500).json({ error: "Error al eliminar el documento" });
            } finally {
                await db.close();
            }
        });

        this.app.listen(this.port, () => {
            new Debug('success').print(`🚀 Server running on http://localhost:${this.port}`);
        });
    }
};