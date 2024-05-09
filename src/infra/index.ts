import { config } from "dotenv";
config();

import API from "../app/handler";

const api = new API();

api.start();