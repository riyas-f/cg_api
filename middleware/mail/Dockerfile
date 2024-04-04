FROM docker.io/library/node:16-alpine@sha256:a1f9d027912b58a7c75be7716c97cfbc6d3099f3a97ed84aa490be9dee20e787


COPY package.json /tmp/package.json
RUN cd /tmp && npm install
RUN mkdir -p /app && cp -a /tmp/node_modules app/

WORKDIR /app

COPY . . 

RUN chmod +x scripts/compose-entrypoint.sh

CMD ["node", "app.js"]