exports = async function (changeEvent) {
    if (!changeEvent.fullDocument) {
        console.error("No fullDocument in changeEvent");
        return;
    }

    const axios = require("axios").default;

    const data = {
        platform: 'mongodb',
        operate: changeEvent.operationType,
        db: changeEvent.ns.db,
        table: changeEvent.ns.coll,
        data: changeEvent.fullDocument
    };

    const config = {
        headers: {
            'Content-Type': 'application/json',
            'Authorization': '*******'
        }
    };

    axios.post('https://example.com/api/event', data, config)
        .then(response => {
            console.log('Response data:', response.data);
        })
        .catch(error => {
            console.error('Error:', error);
        });
};

