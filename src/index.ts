import { config } from 'dotenv';
config();

import express, { Request, Response } from "express";
import bodyParser from "body-parser";
import MongoDB from './MongoDB/MongoConnection';
import { ObjectId } from 'mongodb';
import { color } from './functions';

const app = express();
const port = process.env.PORT || 3000;

app.use(bodyParser.json());

//Busqueda general

app.get("/api/:database/:collection", async (req: Request, res: Response) => {
    const databaseName = req.params.database;
    const mongoDB = new MongoDB(databaseName);

    try {
        await mongoDB.connect();
        const collectionName = req.params.collection;
        const collection = mongoDB.getDatabase().collection(collectionName);
        const documents = await collection.find({}).toArray();
        res.json(documents);
    } catch (error) {
        res.status(500).json({ error: "Error al obtener documentos" });
    } finally {
        await mongoDB.close();
    }
});

// Insercion nueva

app.post("/api/:database/:collection", async (req: Request, res: Response) => {
    const databaseName = req.params.database;
    const mongoDB = new MongoDB(databaseName);

    try {
        await mongoDB.connect();
        const collectionName = req.params.collection;
        const collection = mongoDB.getDatabase().collection(collectionName);
        const newDocument = req.body;

        const result = await collection.insertOne(newDocument);

        if (result) {
            res.status(200).json({
                message: "Data insertada correctamente"
            });
        } else {
            res.status(500).json({ error: "Error al crear el documento" });
        }
    } catch (error) {
        console.error("Error al crear el documento:", error);
        res.status(500).json({ error: "Error al crear el documento" });
    } finally {
        await mongoDB.close();
    }
});

// Actualizacion en base a ID

app.put("/api/:database/:collection/:id", async (req: Request, res: Response) => {
    const databaseName = req.params.database;
    const mongoDB = new MongoDB(databaseName);

    try {
        await mongoDB.connect();
        const collectionName = req.params.collection;
        const collection = mongoDB.getDatabase().collection(collectionName);
        
        const documentId = new ObjectId(req.params.id);

        const updatedDocument = req.body;

        const result = await collection.updateOne({ _id: documentId }, { $set: updatedDocument });
        if (result.modifiedCount === 1) {
            res.json({ message: "Documento actualizado exitosamente" });
        } else {
            res.status(404).json({ error: "Documento no encontrado" });
        }
    } catch (error) {
        res.status(500).json({ error: "Error al actualizar el documento" });
    } finally {
        await mongoDB.close();
    }
});

// Eliminacion en base a ID

app.delete("/api/:database/:collection/:id", async (req: Request, res: Response) => {
    const databaseName = req.params.database;
    const mongoDB = new MongoDB(databaseName);

    try {
        await mongoDB.connect();
        const collectionName = req.params.collection;
        const collection = mongoDB.getDatabase().collection(collectionName);
        
        const documentId = new ObjectId(req.params.id);
        
        const result = await collection.deleteOne({ _id: documentId });

        if (result.deletedCount === 1) {
            res.json({ message: "Documento eliminado exitosamente" });
        } else {
            res.status(404).json({ error: "Documento no encontrado" });
        }
    } catch (error) {
        res.status(500).json({ error: "Error al eliminar el documento" });
    } finally {
        await mongoDB.close();
    }
});

process.on("SIGINT", async () => {
    process.exit(0);
});

app.listen(port, () => {
    // console.log(`API is running on http://localhost:${port}`);
    // const info = 

    console.table([{url: 'http://localhost:', puerto: `${port}`, status: `ON`}])
});