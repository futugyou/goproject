exports = async function (changeEvent) {
    const axios = require("axios").default;

    const data = {
        platform: 'mongodb',
        operate: changeEvent.operationType,
        db: changeEvent.ns.db,
        table: changeEvent.ns.coll,
        data: changeEvent.fullDocument
    };

    if (data.data == undefined) {
        data.data = changeEvent.documentKey
    }

    const config = {
        headers: {
            'Content-Type': 'application/json',
            'Authorization': '*******'
        }
    };

    axios.post('https://example.com/api/v1/event', data, config)
        .then(response => {
            console.log('Response data:', response.data);
        })
        .catch(error => {
            console.error('Error:', error);
        });
};

