import express from 'express';
import routes from './routes';
import middleware from './routes/middleware';

const app = express();

app.use(express.json());
app.use('/', middleware.logger, routes);
app.use(middleware.errorHandler);

app.listen(9080, () => {
    console.log('Listening on port:', 9080);
});
