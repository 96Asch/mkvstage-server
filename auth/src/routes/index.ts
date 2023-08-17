import express from 'express';
import tokenRoute from './token';
import userRoute from './user';

const routes = express.Router();

routes.use('/tokens', tokenRoute);
routes.use('/users', userRoute);

export default routes;
