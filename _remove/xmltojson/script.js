'use strict';

var Jsonix = require('jsonix').Jsonix;
// For instance, in node.js:
var SAM = require('./SAM').SAM;
const fs = require('fs');
var util = require('util');
const json = require('canonicaljson');

// First we construct a Jsonix context - a factory for unmarshaller (parser)
// and marshaller (serializer)
var context = new Jsonix.Context([SAM]);

// Then we create a unmarshaller
var unmarshaller = context.createUnmarshaller();

// Unmarshal an object from the XML retrieved from the URL
unmarshaller.unmarshalFile('./VMP.xml',
    // This callback function will be provided with the result of the unmarshalling
    function (unmarshalled) {
        fs.writeFileSync('/tmp/VMP.json', json.stringify(unmarshalled) , 'utf-8');
    });