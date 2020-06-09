'use strict';

const yargs = require('yargs');
const Xsd2JsonSchema = require('xsd2jsonschema').Xsd2JsonSchema;
const fs = require('fs')
const path = require("path");

const argv = yargs
    .command('convert', 'convert xsd 2 json schema', {
        xsdpath: {
            description: 'path of xsd',
            alias: 'xp',
        },
        jsonpath: {
            description: 'output of json',
            alias: 'jp', 
        }
    })
    .help()
    .alias('help', 'h')
    .argv;

if (argv._.includes('convert')) {
    const xsdPath = argv.xsdpath;
    const xsdName = path.basename(xsdPath);
    const jsonName = path.basename(xsdPath, '.xsd') + '.json';
    fs.access(xsdPath, fs.F_OK, (err) => {
        if (err) {
          console.error(err);
          return;
        }

        console.log(xsdPath)
        fs.readFile(xsdPath, (err, data) => {
            if (err) throw err;
            
            const xs2js = new Xsd2JsonSchema();
            
            const convertedSchemas = xs2js.processAllSchemas({
                schemas: {xsdName: data.toString()}
            });

            const jsonSchema = convertedSchemas[xsdName].get();
            console.log(JSON.stringify(jsonSchema, null, 2));
          });
      })
}



