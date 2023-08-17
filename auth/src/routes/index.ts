import express from 'express';
import authRoute from './auth';
import userRoute from './user';

const routes = express.Router();

routes.use('/auth', authRoute);
routes.use('/users', userRoute);

export default routes;
