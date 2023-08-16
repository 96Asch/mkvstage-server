import express from "express";
import authRoute from "./auth";
import userRoute from "./user";

export const routes = express.Router();

routes.use("/auth", authRoute);
routes.use("/user", userRoute);
