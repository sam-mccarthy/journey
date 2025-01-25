import fs from 'node:fs/promises';
import express from 'express';
import { createServer } from 'vite';

const app = express();

const vite = await createServer({
    server: {
        middlewareMode: true,
    },
    appType: 'custom',
});

app.use(vite.middlewares);
app.use('*', async (request, response) => {
    try {
        const url = request.originalUrl;

        /** @type {string} */
        const raw_html = await fs.readFile('./index.html',{encoding: 'utf-8'});
        const template = await vite.transformIndexHtml(url, raw_html);

        /** @type {import('./src/entry-server.ts').render} */
        const render = (await import('./dist/server/entry-server.js')).render;

        const rendered = await render(url);
        const html = template
            .replace('<!--app-head-->', rendered.head ?? '')
            .replace('<!--app-html-->', rendered.html ?? '');

        response.status(200).set('Content-Type', 'text/html').send(html);
    } catch (err) {
        vite?.ssrFixStacktrace(err);
        console.log(err.stack);
        response.status(500).end(err.stack);
    }
});