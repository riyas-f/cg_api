/*
* This module is for crating axios instance
* with custom https agent that will send
* client certificate to the endpoint
*/

const fs = require('fs');
const https = require('https');
const axios = require('axios');


const CERT_FILE = process.env.CERT_FILE_PATH | 'cert/client.crt'
const KEY_FILE = process.env.PRIVATE_KEY_PATH | 'cert/client.key'
const ROOT_CA = process.env.ROOT_CA_CER | 'cert/root-ca.crt'
const PASSPHRASE_PATH = process.env.PASSPHRASE_PATH | "/tmp/passphrase"

// Read passphrase from file
const passphrase = fs.readFileSync(PASSPHRASE_PATH, 'utf8').trim();

// Load certificate and private key
const cert = fs.readFileSync(CERT_FILE);
const key = fs.readFileSync(KEY_FILE);

// Load the root CA certificate
const caCert = fs.readFileSync(ROOT_CA);

// Create HTTPS agent
const httpsAgent = new https.Agent({
  cert: cert,
  key: key,
  passphrase: passphrase,
  ca: caCert
});

const instance = axios.create({
    httpsAgent: httpsAgent
});

module.exports = {
    axios,
    instance,
}