import fs from 'node:fs/promises';
import express from 'express';
import { createServer } from 'vite';

const options = {
    production: process.env.NODE_ENV === 'production',
    port: process.env.PORT || 8080,
};

const app = express();

const vite = await createServer({
    server: {
        middlewareMode: true,
    },
    appType: 'custom',
});

async function handleRequestProd(url, response) {
    console.log('hai');
    /** @type {string} */
    const template = await fs.readFile('./dist/client/index.html', 'utf8');

    /** @type {import('./dist/server/entry-server.js').render} */
    const render = (await import('./dist/server/entry-server.js')).render;
    const rendered = await render(url);

    return template
        .replace('<!--app-head-->', rendered.head ?? '')
        .replace('<!--app-html-->', rendered.html ?? '');
}

async function handleRequestDev(url, response) {
    /** @type {string} */
    const raw = await fs.readFile('./index.html', 'utf8');
    const template = await vite.transformIndexHtml(url, raw);

    /** @type {import('./src/entry-server.tsx').render} */
    const render = (await vite.ssrLoadModule('/src/entry-server.tsx')).render;
    const rendered = render(url);

    return template
        .replace('<!--app-head-->', rendered.head ?? '')
        .replace('<!--app-html-->', rendered.html ?? '');
}

app.use(vite.middlewares);
app.use('*', async (request, response) => {
    try {
        /** @type {string} */
        let html;
        if(options.production)
            html = await handleRequestProd(request.originalUrl, response);
        else
            html = await handleRequestDev(request.originalUrl, response);
        response.status(200).set('Content-Type', 'text/html').send(html);
    } catch (err) {
        vite?.ssrFixStacktrace(err);
        console.log(err.stack);
        response.status(500).end(err.stack);
    }
});

app.listen(8080, () => {
    console.log(`Server started at http://localhost:8080`);
});